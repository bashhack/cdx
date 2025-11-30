// User handler for API requests

export interface User {
  id: number;
  name: string;
  email: string;
}

export class UserHandler {
  constructor(private repository: UserRepository) {}

  async getUser(id: number): Promise<User | null> {
    return this.repository.findById(id);
  }
}

export function createHandler(repo: UserRepository): UserHandler {
  return new UserHandler(repo);
}

export const fetchUser = async (id: number): Promise<User> => {
  const response = await fetch(`/api/users/${id}`);
  return response.json();
};

interface UserRepository {
  findById(id: number): Promise<User | null>;
}

export type UserId = number;
