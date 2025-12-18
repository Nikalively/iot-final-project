from locust import HttpUser, task, between
import json
import time
import random

class IoTUser(HttpUser):
    wait_time = between(0.001, 0.01)

    @task
    def post_metrics(self):
        current_time = time.strftime('%Y-%m-%dT%H:%M:%S.%fZ', time.gmtime())
        payload = {
            "timestamp": current_time,
            "cpu": random.uniform(20.0, 80.0),
            "rps": random.uniform(5.0, 15.0) if random.random() > 0.1 else random.uniform(50.0, 100.0)
        }
        self.client.post("/metrics", json=payload)
