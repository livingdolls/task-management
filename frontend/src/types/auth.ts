export type TRegister = {
  name: string;
  username: string;
  password: string;
};

export type TLogin = Pick<TRegister, "username" | "password">;

export type TUser = {
  id: number;
  username: string;
  name: string;
};

export type TLoginResponse = {
  token: string;
  user: TUser;
};

export type TRegisterResponse = TUser;
