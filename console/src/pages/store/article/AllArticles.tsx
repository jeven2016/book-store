import {Box, Container, Space, Table} from 'react-windy-ui';
import dayjs from 'dayjs';
import {useCallback, useEffect, useState} from 'react';
import {ArticlePage, Catalog} from '@/Types';
import {get} from '@/client/Request';
import {buildUrl} from '@/common/utils';

export default function AllArticles(props) {
    const [tableData, setTableData] = useState<ArticlePage | null>(null);
    const [catalogs, setCatalogs] = useState<Catalog[]>([]);
    const [subCatalogId, setSubCatalogId] = useState<string | null>(null);
    const [page, setPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(20);

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
            catalogId: subCatalogId
        }).then((data) => {
            setTableData(data as ArticlePage);
        });
    }, [subCatalogId, page, pageSize]);

    useEffect(() => {
        if (!subCatalogId) {
            return;
        }
        showDefaultPage();
    }, [subCatalogId]);

    return (
        <Container extraClassName="bs-content" autoAdjust={true}>
            <div className="bs-content-panel with-gutter">
                {catalogs.map((c: Catalog) => (
                    <Box
                        key={`cat-${c.id}`}
                        block
                        left="武侠传奇："
                        center={
                            <Space gutter={{x: 32}}>
                                {c.children.map((chd: Catalog) => (
                                    <span
                                        className={`bs-link ${subCatalogId === c.id ? 'active' : ''}`}
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
                target="_blank">
                {text}
            </a>
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
