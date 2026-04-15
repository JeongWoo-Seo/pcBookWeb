import ws from 'k6/ws';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 100 },
    { duration: '2m', target: 500 },
    { duration: '3m', target: 1000 },
  ],
};

export default function () {
  const url = 'ws://localhost:8081/ws?mode=single&id=FVFFP6NDQ6L4';

  const res = ws.connect(url, {}, function (socket) {

    socket.on('open', () => {
      // 연결 성공
    });

    socket.on('message', (data) => {
      const msg = JSON.parse(data);
      const createdAtMs = (msg.createAt.seconds * 1000) + (msg.createAt.nanos / 1000000);
      const latency = Date.now() - createdAtMs;

      console.log(`Latency: ${latency}ms`);
    });

    socket.on('error', (e) => {
      console.error('error:', e.error());
    });

    socket.on('close', () => {
      // 연결 끊김 추적 가능
    });

    // ping 유지
    socket.setInterval(() => {
      socket.send(JSON.stringify({ type: 'ping' }));
    }, 5000);

    socket.setTimeout(() => {
      socket.close();
    }, 10000);
  });

  check(res, {
    'connection successful': (r) => r && r.status === 101,
  });

  sleep(1);
}