// Rust user module

pub struct User {
    pub id: i64,
    pub name: String,
    pub email: String,
}

pub trait UserRepository {
    fn find_by_id(&self, id: i64) -> Option<User>;
    fn create(&self, user: &User) -> Result<(), String>;
}

pub struct UserService<R: UserRepository> {
    repository: R,
}

impl<R: UserRepository> UserService<R> {
    pub fn new(repository: R) -> Self {
        UserService { repository }
    }

    pub fn get_user(&self, id: i64) -> Option<User> {
        self.repository.find_by_id(id)
    }
}

pub fn create_user(name: String, email: String) -> User {
    User { id: 0, name, email }
}

pub async fn fetch_user(id: i64) -> Result<User, String> {
    todo!()
}

pub enum UserRole {
    Admin,
    Member,
    Guest,
}
