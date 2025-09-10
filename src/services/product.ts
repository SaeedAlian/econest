import { CancelToken } from "axios";

import api from "@/lib/axios";

import {
  IProduct,
  IProductPageCountResponse,
  IProductSearchQuery,
  IProductExtended,
  IProductCommentSearchQuery,
  IProductCommentWithUser,
  IProductCommentPageCountResponse,
  parseProduct,
  parseProductExtended,
  parseProductCommentWithUser,
  IProductTagSearchQuery,
  IProductTag,
  parseProductTag,
  IProductTagPageCountResponse,
} from "@/types/product";

export class ProductService {
  static baseURL = "http://localhost:5000";

  static async getProducts(
    params: IProductSearchQuery = {},
    cancelToken?: CancelToken | undefined,
  ): Promise<IProduct[]> {
    const res = await api.get<IProduct[]>(`${this.baseURL}/product`, {
      params,
      cancelToken,
    });
    return res.data.map(parseProduct);
  }

  static async getProductPages(
    params: IProductSearchQuery = {},
  ): Promise<number> {
    const res = await api.get<IProductPageCountResponse>(
      `${this.baseURL}/product/pages`,
      {
        params,
      },
    );
    return res.data.pages;
  }

  static async getProductTags(
    params: IProductTagSearchQuery = {},
    cancelToken?: CancelToken | undefined,
  ): Promise<IProductTag[]> {
    const res = await api.get<IProductTag[]>(`${this.baseURL}/product/tag`, {
      params,
      cancelToken,
    });
    return res.data.map(parseProductTag);
  }

  static async getProductTagsPages(
    params: IProductTagSearchQuery = {},
  ): Promise<number> {
    const res = await api.get<IProductTagPageCountResponse>(
      `${this.baseURL}/product/tag/pages`,
      {
        params,
      },
    );
    return res.data.pages;
  }

  static async getProductTag(tagId: number): Promise<IProductTag> {
    const res = await api.get<IProductTag>(
      `${this.baseURL}/product/tag/${tagId}`,
    );
    return parseProductTag(res.data);
  }

  static async getProductComments(
    productId: number,
    params: IProductCommentSearchQuery = {},
    cancelToken?: CancelToken | undefined,
  ): Promise<IProductCommentWithUser[]> {
    const res = await api.get<IProductCommentWithUser[]>(
      `${this.baseURL}/product/comment/withuser/product/${productId}`,
      {
        params,
        cancelToken,
      },
    );
    return res.data.map(parseProductCommentWithUser);
  }

  static async getProductCommentsPages(
    productId: number,
    params: IProductCommentSearchQuery = {},
  ): Promise<number> {
    const res = await api.get<IProductCommentPageCountResponse>(
      `${this.baseURL}/product/comment/product/${productId}/pages`,
      {
        params,
      },
    );
    return res.data.pages;
  }

  static async getProduct(productId: number): Promise<IProduct> {
    const res = await api.get<IProduct>(`${this.baseURL}/product/${productId}`);
    return parseProduct(res.data);
  }

  static async getProductExtended(
    productId: number,
  ): Promise<IProductExtended> {
    const res = await api.get<IProductExtended>(
      `${this.baseURL}/product/${productId}/extended`,
    );
    return parseProductExtended(res.data);
  }
}
