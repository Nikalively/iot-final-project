from locust import HttpUser, between, task
import time
import random

class MetricsUser(HttpUser):
    wait_time = between(0.005, 0.02)

    @task
    def post_metrics(self):
        timestamp = int(time.time())

        cpu = round(random.uniform(0, 100), 2)
        rps = round(random.uniform(0, 100), 2)

        anomaly = random.random() < 0.1
        if anomaly:
            rps = round(random.uniform(50, 100), 2)

        payload = {
            'timestamp': timestamp,
            'cpu_usage': cpu,
            'rps': rps
        }

        with self.client.post("/metrics", json=payload, catch_response=True) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"HTTP {response.status_code}: {response.text}")

