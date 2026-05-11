import socket
import struct
import msgpack
import random

def send_msg(sock, msg_id, data):
    data_len = len(data)
    # 调整顺序：MsgID 在前 (I)，DataLen 在后 (I)
    header = struct.pack('>II', msg_id, data_len)
    sock.sendall(header + data)

def recv_msg(sock):
    header_data = sock.recv(8)
    if not header_data:
        return None
    # 接收时也按相同顺序解析
    msg_id, data_len = struct.unpack('>II', header_data)
    body = sock.recv(data_len)
    return msg_id, body


def main():
    server_addr = ('127.0.0.1', 8888)
    
    try:
        # 创建 TCP Socket
        client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client.connect(server_addr)
        
        # 发送测试数据 (MsgID=0, 跟 main.go 里的 AddRouter 对应)
        test_sum = msgpack.packb({
            "arg1": random.randint(1, 99999),
            "arg2": random.randint(1,99999),
            "result": 0,
        })
        print(f"[Send] MsgID: 1001, Data: {test_sum}")
        send_msg(client, 1001, test_sum)
        
        # 接收服务器回执
        res = recv_msg(client)
        if res:
            mid, data = res
            print(f"[Recv] MsgID: {mid}, Data: {msgpack.unpackb(data)}")
            
    except Exception as e:
        print(f"Error: {e}")
    finally:
        client.close()

if __name__ == "__main__":
    main()
