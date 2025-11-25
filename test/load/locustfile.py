from locust import HttpUser, task, between
import uuid
import random

# 20 команд * 10 пользователей = 200 author_id
AUTHORS = [
    f"user-loadtest-{team:02d}-{member:02d}"
    for team in range(1, 21)
    for member in range(1, 11)
]

class PRReviewerUser(HttpUser):
    host = "http://localhost:8080"
    wait_time = between(1, 3)

    def on_start(self):
        self.author_id = random.choice(AUTHORS)
        self._create_initial_pr()

    def _create_initial_pr(self):
        pr_id = f"pr-{uuid.uuid4().hex}"
        payload = {
            "pull_request_id": pr_id,
            "pull_request_name": "Initial LoadTest PR",
            "author_id": self.author_id,
        }

        response = self.client.post(
            "/pullRequest/create",
            json=payload,
            name="/pullRequest/create",
        )

        if response.status_code == 201:
            data = response.json()
            self.current_pr_id = pr_id

            pr_data = data.get("pr", {})
            assigned = pr_data.get("assigned_reviewers", [])
            if assigned:
                self.current_reviewer_to_replace = assigned[0]

    @task(2)
    def create_pull_request(self):
        pr_id = f"pr-{uuid.uuid4().hex}"
        payload = {
            "pull_request_id": pr_id,
            "pull_request_name": "LoadTest PR",
            "author_id": self.author_id,
        }

        response = self.client.post(
            "/pullRequest/create",
            json=payload,
            name="/pullRequest/create",
        )

        if response.status_code == 201:
            data = response.json()
            self.current_pr_id = pr_id

            pr_data = data.get("pr", {})
            assigned = pr_data.get("assigned_reviewers", [])
            if assigned:
                self.current_reviewer_to_replace = assigned[0]

    @task(1)
    def reassign_reviewer(self):
        if not hasattr(self, "current_pr_id") or not hasattr(self, "current_reviewer_to_replace"):
            return

        payload = {
            "pull_request_id": self.current_pr_id,
            "old_user_id": self.current_reviewer_to_replace,
        }

        response = self.client.post(
            "/pullRequest/reassign",
            json=payload,
            name="/pullRequest/reassign",
        )

        if response.status_code == 200:
            data = response.json()
            pr_data = data.get("pr", {})
            assigned = pr_data.get("assigned_reviewers", [])
            if assigned:
                self.current_reviewer_to_replace = assigned[0]
