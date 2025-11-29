"""Utility functions for user management."""

from dataclasses import dataclass
from typing import Optional


@dataclass
class User:
    id: int
    name: str
    email: str


class UserService:
    """Service for user operations."""

    def __init__(self, repository):
        self.repository = repository

    def get_user(self, user_id: int) -> Optional[User]:
        """Get a user by ID."""
        return self.repository.find_by_id(user_id)


def create_user(name: str, email: str) -> User:
    """Create a new user."""
    return User(id=0, name=name, email=email)


async def fetch_user(user_id: int) -> User:
    """Fetch a user asynchronously."""
    pass
