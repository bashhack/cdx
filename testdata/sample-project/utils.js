// JavaScript utility functions

class UserManager {
  constructor(repository) {
    this.repository = repository;
  }

  async getUser(id) {
    return this.repository.findById(id);
  }
}

function createUser(name, email) {
  return { id: 0, name, email };
}

const fetchUserData = async (userId) => {
  const response = await fetch(`/api/users/${userId}`);
  return response.json();
};

async function deleteUser(id) {
  await fetch(`/api/users/${id}`, { method: 'DELETE' });
}

module.exports = { UserManager, createUser, fetchUserData, deleteUser };
