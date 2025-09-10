import { parseDate } from "@/lib/utils";
import { ICommentUser, parseCommentUser } from "@/types/user";

export interface IProductBase {
  id: number;
  name: string;
  slug: string;
  price: number;
  shipmentFactor: number;
  description: string;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
  subcategoryId: number;
}

export interface IProductOffer {
  id: number;
  discount: number;
  expireAt: Date;
  createdAt: Date;
  updatedAt: Date;
  productId: number;
}

export interface IProductCategory {
  id: number;
  name: string;
  imageName: string;
  createdAt: Date;
  updatedAt: Date;
  parentCategoryId?: number;
}

export interface IProductCategoryWithParents extends IProductCategory {
  parentCategory?: IProductCategoryWithParents | null;
}

export interface IProductTag {
  id: number;
  name: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface IProductSpec {
  id: number;
  label: string;
  value: string;
  productId: number;
}

export interface IProductTagAssignment {
  productId: number;
  tagId: number;
}

export interface IProductAttribute {
  id: number;
  label: string;
}

export interface IProductAttributeOption {
  id: number;
  value: string;
  attributeId: number;
}

export interface IProductAttributeWithOptions extends IProductAttribute {
  options: IProductAttributeOption[];
}

export interface IProductVariant {
  id: number;
  quantity: number;
  productId: number;
}

export interface IProductVariantAttributeOption {
  variantId: number;
  attributeId: number;
  optionId: number;
}

export interface IProductVariantSelectedAttributeOption
  extends IProductAttribute {
  selectedOption: IProductAttributeOption;
}

export interface IProductVariantWithAttributeSet extends IProductVariant {
  attributeSet: IProductVariantSelectedAttributeOption[];
}

export interface IProductComment {
  id: number;
  scoring: number;
  comment?: string;
  createdAt: Date;
  updatedAt: Date;
  productId: number;
  userId: number;
}

export interface IProductCommentWithUser {
  id: number;
  scoring: number;
  comment?: string;
  createdAt: Date;
  updatedAt: Date;
  productId: number;
  user: ICommentUser;
}

export interface IProductImage {
  id: number;
  imageName: string;
  isMain: boolean;
  productId: number;
}

export interface IStoreInfo {
  id: number;
  name: string;
  description: string;
}

export interface IProduct extends IProductBase {
  subcategory: IProductCategory;
  averageScore: number;
  totalQuantity: number;
  offer?: IProductOffer;
  mainImage?: IProductImage;
  store: IStoreInfo;
}

export interface IProductExtended extends IProductBase {
  subcategory: IProductCategoryWithParents;
  specs: IProductSpec[];
  tags: IProductTag[];
  variants: IProductVariantWithAttributeSet[];
  attributes: IProductAttributeWithOptions[];
  offer?: IProductOffer | null;
  images: IProductImage[];
  store: IStoreInfo;
}

export interface IProductSearchQuery {
  k?: string;
  avgscr?: number;
  minq?: number;
  maxq?: number;
  offr?: boolean;
  cat?: number;
  tags?: string;
  pmt?: number;
  plt?: number;
  store?: number;
  p?: number;
  offst?: number;
  lim?: number;
}

export interface IProductTagSearchQuery {
  name?: string;
  product?: string;
  p?: number;
  offst?: number;
  lim?: number;
}

export interface IProductCommentSearchQuery {
  slt?: number;
  smt?: number;
  p?: number;
}

export interface IProductPageCountResponse {
  pages: number;
}

export interface IProductTagPageCountResponse {
  pages: number;
}

export interface IProductCommentPageCountResponse {
  pages: number;
}

export function parseProductOffer(raw: any): IProductOffer {
  return {
    ...raw,
    expireAt: parseDate(raw.expireAt)!,
    createdAt: parseDate(raw.createdAt)!,
    updatedAt: parseDate(raw.updatedAt)!,
  };
}

export function parseProductCategory(raw: any): IProductCategory {
  return {
    ...raw,
    createdAt: parseDate(raw.createdAt)!,
    updatedAt: parseDate(raw.updatedAt)!,
  };
}

export function parseProductCategoryWithParents(
  raw: any,
): IProductCategoryWithParents {
  return {
    ...parseProductCategory(raw),
    parentCategory: raw.parentCategory
      ? parseProductCategoryWithParents(raw.parentCategory)
      : undefined,
  };
}

export function parseProductTag(raw: any): IProductTag {
  return {
    ...raw,
    createdAt: parseDate(raw.createdAt)!,
    updatedAt: parseDate(raw.updatedAt)!,
  };
}

export function parseProductComment(raw: any): IProductComment {
  return {
    ...raw,
    createdAt: parseDate(raw.createdAt)!,
    updatedAt: parseDate(raw.updatedAt)!,
  };
}

export function parseProductCommentWithUser(raw: any): IProductCommentWithUser {
  return {
    ...raw,
    createdAt: parseDate(raw.createdAt)!,
    updatedAt: parseDate(raw.updatedAt)!,
    user: parseCommentUser(raw.user),
  };
}

export function parseProductBase(raw: any): IProductBase {
  return {
    ...raw,
    createdAt: parseDate(raw.createdAt)!,
    updatedAt: parseDate(raw.updatedAt)!,
  };
}

export function parseProduct(raw: any): IProduct {
  return {
    ...raw,
    createdAt: parseDate(raw.createdAt)!,
    updatedAt: parseDate(raw.updatedAt)!,
    offer: raw.offer ? parseProductOffer(raw.offer) : null,
  };
}

export function parseProductExtended(raw: any): IProductExtended {
  return {
    ...raw,
    createdAt: parseDate(raw.createdAt)!,
    updatedAt: parseDate(raw.updatedAt)!,
    subcategory: parseProductCategoryWithParents(raw.subcategory),
    specs: raw.specs ?? [],
    tags: raw.tags ? raw.tags.map(parseProductTag) : [],
    variants: raw.variants ?? [],
    offer: raw.offer ? parseProductOffer(raw.offer) : null,
    images: raw.images ?? [],
    store: raw.store,
  };
}
