package bookstore.queryapi.mapper;

import bookstore.queryapi.document.ArticleCatalogDoc;
import bookstore.queryapi.dto.ArticleCatalogDto;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

@Mapper
public interface ArticleCatalogMapper {
    ArticleCatalogMapper INSTANCE = Mappers.getMapper(ArticleCatalogMapper.class);

    ArticleCatalogDto toDto(ArticleCatalogDoc doc);

    ArticleCatalogDoc toDocument(ArticleCatalogDto doc);
}
