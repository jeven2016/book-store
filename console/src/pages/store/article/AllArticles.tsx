import { Box, Container, Pagination, Space, Table } from 'react-windy-ui';
import dayjs from 'dayjs';
import { useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { ArticlePage, ArticleSearchContextValue, Catalog, WindowChangeInfo } from '@/Types';
import { get } from '@/client/Request';
import { buildUrl } from '@/common/utils';
import { useTranslation } from 'react-i18next';
import { ArticleSearchCtx, WindowChangeContext } from '@/common/Context';

export default function AllArticles(props) {
  const { t } = useTranslation();
  const [tableData, setTableData] = useState<ArticlePage | null>(null);
  const [catalogs, setCatalogs] = useState<Catalog[]>([]);
  const [subCatalogId, setSubCatalogId] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(20);
  const { sm }: WindowChangeInfo = useContext<WindowChangeInfo>(WindowChangeContext);

  const { name } = useContext<ArticleSearchContextValue>(ArticleSearchCtx);

  useEffect(() => {
    get(buildUrl('/article-catalogs')).then((data) => {
      const list = data as Catalog[];
      setCatalogs(list);
    });
  }, []);

  const changeCatalog = useCallback((id) => {
    setSubCatalogId(id);
    setPage(1);
  }, []);

  const showDefaultPage = useCallback(() => {
    const pageUrl = `/articles?page=${page}&limit=${pageSize}`;
    get(buildUrl(pageUrl), {
      page,
      limit: pageSize,
      catalogId: subCatalogId,
      name
    }).then((data) => {
      setTableData(data as ArticlePage);
    });
  }, [subCatalogId, page, pageSize, name]);

  useEffect(() => {
    showDefaultPage();
  }, [showDefaultPage]);

  const goTo = useCallback((page: number, limit: number) => {
    setPageSize(limit);
    setPage(page);
  }, []);

  const changePageSize = useCallback((pageSize: number) => {
    setPageSize(pageSize);
  }, []);

  const cells = useMemo(() => (sm ? getSmCells(catalogs) : getCells(catalogs)), [sm, catalogs]);

  return (
    <Container extraClassName="bs-content" autoAdjust={true}>
      <div className="bs-content-panel with-gutter">
        {catalogs.map((c: Catalog) => (
          <div key={`cat-${c.id}`}>
            <Box
              block
              extraClassName="bs-all-catalogs"
              left={c.name + ':'}
              center={
                <Space gutter={{ x: 16, y: 8 }} wrap>
                  <span
                    className={`bs-link ${!subCatalogId ? 'active' : ''}`}
                    onClick={() => changeCatalog(null)}>
                    不限
                  </span>
                  {c.children.map((chd: Catalog) => (
                    <span
                      className={`bs-link ${subCatalogId === chd.id ? 'active' : ''}`}
                      key={`chd-${chd.id}`}
                      onClick={() => changeCatalog(chd.id)}>
                      {chd.name}
                    </span>
                  ))}
                </Space>
              }
            />
          </div>
        ))}
      </div>
      <div className="bs-content-panel with-gutter bs-book-list">
        <Table
          type="simple"
          hover={true}
          loadData={tableData?.rows || []}
          cells={cells}
          checkable={false}
          checkType="checkbox"
          onCheckAll={(next) => console.log('check all: ' + next)}
          onCheckChange={(jsonData, next) => console.log('check one: ' + jsonData + next)}
        />
        <Box
          block
          align="end"
          left={<></>}
          right={
            <Pagination
              simple={sm}
              pageCount={tableData?.totalPages}
              defaultPage={1}
              siblingCount={1}
              hasPageRange={true}
              pageRange={pageSize}
              onChangeRange={changePageSize}
              pageRanges={[10, 20, 50, 100]}
              leftItems={[
                `${t('global.pagination.total')}${tableData?.totalPages || 0}${t(
                  'global.pagination.pages'
                )}， ${tableData?.count || 0}${t('global.pagination.records')}`
              ]}
              onChange={goTo}
            />
          }
        />
      </div>
    </Container>
  );
}

const getSmCells = (catalogs: Catalog[]) => {
  return [
    {
      head: '序号',
      paramName: 'key',
      width: '50px',
      format: (text, row, tableIndex) => {
        return (
          <h5>
            <span className="bs-tbl-index">{tableIndex + 1}</span>
          </h5>
        );
      }
    },
    {
      head: '搜索结果',
      paramName: 'name',
      format: (text, row) => (
        <div>
          <a
            className="bs-article-name-link"
            rel="noreferrer"
            href={`/articles/${row.id}`}
            target="_blank"
            dangerouslySetInnerHTML={{ __html: text }}></a>
          <h5>
            <span>[{genCellData(catalogs, row.catalogId)}]</span>
          </h5>
        </div>
      )
    }
  ];
};

const getCells = (catalogs: Catalog[]) => {
  return [
    {
      head: '序号',
      paramName: 'key',
      width: '50px',
      format: (text, row, tableIndex) => {
        return (
          <h5>
            <span className="bs-tbl-index">{tableIndex + 1}</span>
          </h5>
        );
      }
    },
    {
      head: '频道',
      paramName: 'catalogId',
      width: '150px',
      format: (text, row, tableIndex) => {
        return (
          <h5>
            <span>[{genCellData(catalogs, row.catalogId)}]</span>
          </h5>
        );
      }
    },
    {
      head: '书名',
      paramName: 'name',
      format: (text, row) => (
        <a
          className="bs-article-name-link"
          rel="noreferrer"
          href={`/articles/${row.id}`}
          target="_blank"
          dangerouslySetInnerHTML={{ __html: text }}></a>
      )
    },
    {
      head: '入库时间',
      paramName: 'createDate',
      width: '200px',
      format: (text) => {
        return <h5>{dayjs(text).format('YYYY-MM-DD')}</h5>;
      }
    }
  ];
};

const genCellData = (catalogs: Catalog[], catalogId) => {
  for (const c of catalogs) {
    if (catalogId === c.id) {
      return c.name;
    }
    const realName = genCellData(c.children, catalogId);
    if (realName) {
      return realName;
    }
  }
  return null;
};
