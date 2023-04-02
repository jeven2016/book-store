import { AxiosResponse } from 'axios';

interface Tokens {
  idToken?: string;
  refreshToken?: string;
  token?: string;
}

interface WindowContextValue {
  mdWindow?: boolean;
}

interface ResponseData extends AxiosResponse {
  data: { errCode: number };
}

interface Catalog {
  id: string;
  name: string;
  catalogId: string;
  order: number;
  children: Catalog[];
}

interface ArticleCatalogContextValue {
  catalogs: Catalog[];
  selectedCatalogId: string | null;
  changeArticleCatalog: ((id: string) => void) | null;
}

interface ArticleSearchContextValue {
  name: string;
}

interface ArticlePage {
  count: number;
  limit: number;
  page: number;
  totalPages: number;
  rows: Article[];
}

interface ArticleInfo {
  id: string;
  name: string;
  content: string;
  createDate: Date;
}
