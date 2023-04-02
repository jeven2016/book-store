import { Box, Container, Pagination, Space, Table } from 'react-windy-ui';
import dayjs from 'dayjs';
import { useCallback, useContext, useEffect, useState } from 'react';
import { ArticlePage, ArticleSearchContextValue, Catalog } from '@/Types';
import { get } from '@/client/Request';
import { buildUrl } from '@/common/utils';
import { useTranslation } from 'react-i18next';
import { ArticleSearchCtx } from '@/common/Context';

export default function AllArticles(props) {
  const { t } = useTranslation();
  const [tableData, setTableData] = useState<ArticlePage | null>(null);
  const [catalogs, setCatalogs] = useState<Catalog[]>([]);
  const [subCatalogId, setSubCatalogId] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(20);

  const { name } = useContext<ArticleSearchContextValue>(ArticleSearchCtx);

  console.log('name=' + name);

  useEffect(() => {
    get(buildUrl('/article-catalogs')).then((data) => {
      const list = data as Catalog[];
      setCatalogs(list);
    });
  }, []);

  const changeCatalog = useCallback((id) => {
    setSubCatalogId(id);
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

  return (
    <Container extraClassName="bs-content" autoAdjust={true}>
      <div className="bs-content-panel with-gutter">
        {catalogs.map((c: Catalog) => (
          <Box
            key={`cat-${c.id}`}
            block
            left={c.name + ':'}
            center={
              <Space gutter={{ x: 32 }}>
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
        ))}
      </div>
      <div className="bs-content-panel with-gutter bs-book-list">
        <Box
          block
          align="end"
          left={<></>}
          right={
            <Pagination
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

const cells = [
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
