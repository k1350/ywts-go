import AuthRepository from "./authRepository";
import BoardsRepository from "./boardsRepository";

const repositories = {
  auth: AuthRepository,
  boards: BoardsRepository,
};

export const RepositoryFactory = {
  get: (name) => repositories[name],
};
