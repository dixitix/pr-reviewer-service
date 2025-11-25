#!/usr/bin/env python3
import requests

BASE_URL = "http://localhost:8080"

NUM_TEAMS = 20
USERS_PER_TEAM = 10

def main():
    for team_index in range(1, NUM_TEAMS + 1):
        team_name = f"loadtest-team-{team_index:02d}"

        members = []
        for user_index in range(1, USERS_PER_TEAM + 1):
            user_id = f"user-loadtest-{team_index:02d}-{user_index:02d}"
            username = f"loadtest-user-{team_index:02d}-{user_index:02d}"
            members.append({
                "user_id": user_id,
                "username": username,
                "is_active": True,
            })

        payload = {
            "team_name": team_name,
            "members": members,
        }

        resp = requests.post(f"{BASE_URL}/team/add", json=payload, timeout=5)
        print(f"[team {team_name}] status={resp.status_code}")
        if resp.status_code not in (200, 201, 409):
            print("  body:", resp.text)
            resp.raise_for_status()

if __name__ == "__main__":
    main()
