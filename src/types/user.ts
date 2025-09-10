import { parseDate } from "@/lib/utils";

export interface ICommentUser {
  id: number;
  fullName?: string;
  createdAt: Date;
  updatedAt: Date;
}

export function parseCommentUser(raw: any): ICommentUser {
  return {
    ...raw,
    createdAt: parseDate(raw.createdAt)!,
    updatedAt: parseDate(raw.updatedAt)!,
  };
}
