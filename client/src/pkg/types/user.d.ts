export interface IUser {
    id: number
    role: "root" | "admin" | "user"
    email: string
    password?: string | null
}