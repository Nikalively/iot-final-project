from locust import HttpUser, task, between
import json
import time

class IoTUser(HttpUser):
    wait_time = between(0.1, 1)

    @task
    def post_metrics(self):
        payload = {
            "timestamp": time.time(),
            "cpu": 50.0,
            "rps": 10.0
        }
        self.client.post("/metrics", json=payload)
