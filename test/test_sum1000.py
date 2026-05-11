import asyncio
import struct
import msgpack
import random
import time

# 统计数据
class Stats:
    def __init__(self):
        self.success = 0
        self.failed = 0

stats = Stats()
# 并发限制：防止瞬间创建1000个TCP连接导致系统拒绝或丢包
sem = asyncio.Semaphore(200) 

async def send_msg(writer, msg_id, data):
    data_len = len(data)
    header = struct.pack('>II', msg_id, data_len)
    writer.write(header + data)
    await writer.drain()

async def recv_msg(reader):
    try:
        # readexactly 保证读够 8 字节，读不够会抛出异常
        header_data = await reader.readexactly(8)
        msg_id, data_len = struct.unpack('>II', header_data)
        body = await reader.readexactly(data_len)
        return msg_id, body
    except:
        return None

async def run_client(client_id, host, port, repeat):
    async with sem: # 使用信号量平滑建立连接
        try:
            reader, writer = await asyncio.open_connection(host, port)
            
            for _ in range(repeat):
                test_data = msgpack.packb({
                    "arg1": random.randint(1, 99999),
                    "arg2": random.randint(1, 99999),
                    "result": 0,
                })
                
                try:
                    await send_msg(writer, 1001, test_data)
                    res = await recv_msg(reader)
                    
                    if res:
                        stats.success += 1
                    else:
                        stats.failed += 1
                except:
                    stats.failed += 1
            
            writer.close()
            await writer.wait_closed()
        except Exception:
            stats.failed += repeat # 建连失败，该连接的任务全部计入失败

async def main():
    server_host = '127.0.0.1'
    server_port = 8888
    total_conns = 100
    repeat_per_conn = 100 
    
    print(f"🚀 启动压测: {total_conns} 并发连接, 每个连接请求 {repeat_per_conn} 次")
    start_time = time.time()

    tasks = [run_client(i, server_host, server_port, repeat_per_conn) for i in range(total_conns)]
    await asyncio.gather(*tasks)

    duration = time.time() - start_time

    # 修正 QPS 计算：仅统计成功的请求
    qps = stats.success / duration if duration > 0 else 0

    print("\n" + "="*40)
    print(f"🏁 压测报告")
    print(f"总耗时: {duration:.2f} 秒")
    print(f"成功次数 (分子): {stats.success}")
    print(f"失败次数: {stats.failed}")
    print(f"有效 QPS (成功/时间): {qps:.2f}")
    print("="*40)

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        pass
