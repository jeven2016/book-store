import React from 'react';
import {
  ArticleCatalogContextValue,
  ArticleSearchContextValue,
  WindowChangeInfo,
  WindowContextValue
} from '@/Types';

export const WindowContext: React.Context<WindowContextValue> =
  React.createContext<WindowContextValue>({ mdWindow: false });

export const ArticleCatalogContext: React.Context<ArticleCatalogContextValue> =
  React.createContext<ArticleCatalogContextValue>({
    catalogs: [],
    selectedCatalogId: null,
    changeArticleCatalog: null
  });

export const WindowChangeContext: React.Context<WindowChangeInfo> =
  React.createContext<WindowChangeInfo>({
    sm: false
  });

export const ArticleSearchCtx: React.Context<ArticleSearchContextValue> =
  React.createContext<ArticleSearchContextValue>({ name: '' });
