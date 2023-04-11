import { Box, Pagination, Table } from 'react-windy-ui';
import { useCallback, useEffect, useState } from 'react';
import { buildUrl } from '@/common/utils';
import { get } from '@/client/Request';
import { ArticlePage } from '@/Types';
import dayjs from 'dayjs';
import { useTranslation } from 'react-i18next';

export default function ArticleList(props) {
  const { subCatalogId, smWindow = false } = props;
  const { t } = useTranslation();
  const [tableData, setTableData] = useState<ArticlePage | null>(null);
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(20);

  useEffect(() => {
    if (!subCatalogId) {
      return;
    }
    const pageUrl = `/article-catalogs/${subCatalogId}/articles?page=${page}&limit=${pageSize}`;
    get(buildUrl(pageUrl)).then((data) => {
      setTableData(data as ArticlePage);
    });
  }, [subCatalogId, page, pageSize]);

  const goTo = useCallback((page: number, limit: number) => {
    setPageSize(limit);
    setPage(page);
  }, []);

  const changePageSize = useCallback((pageSize: number) => {
    setPageSize(pageSize);
  }, []);
  return (
    <div className="bs-book-list">
      <Box
        block
        align="end"
        left={<></>}
        right={
          <Pagination
            simple={smWindow}
            pageCount={tableData?.totalPages}
            page={page}
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
            simple={smWindow}
            pageCount={tableData?.totalPages}
            page={page}
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
  );
}

const cells = [
  {
    head: '序号',
    paramName: 'key',
    width: '50px',
    format: (text, row, tableIndex) => {
      return tableIndex + 1;
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
        target="_blank">
        {text}
      </a>
    )
  },
  {
    head: '入库时间',
    paramName: 'createDate',
    width: '130px',
    format: (text) => {
      return dayjs(text).format('YY-MM-DD');
    }
  }
];
