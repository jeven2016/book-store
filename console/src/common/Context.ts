import React from 'react';
import { ArticleCatalogContextValue, ArticleSearchContextValue, WindowContextValue } from '@/Types';

export const WindowContext: React.Context<WindowContextValue> =
  React.createContext<WindowContextValue>({ mdWindow: false });

export const ArticleCatalogContext: React.Context<ArticleCatalogContextValue> =
  React.createContext<ArticleCatalogContextValue>({
    catalogs: [],
    selectedCatalogId: null,
    changeArticleCatalog: null
  });

export const ArticleSearchCtx: React.Context<ArticleSearchContextValue> =
  React.createContext<ArticleSearchContextValue>({ name: '' });
