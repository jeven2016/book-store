package bookstore.queryapi.mapper;

import bookstore.queryapi.document.ArticleCatalogDoc;
import bookstore.queryapi.dto.ArticleCatalogDto;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

@Mapper
public interface ArticleCatalogMapper {
    ArticleCatalogMapper INSTANCE = Mappers.getMapper(ArticleCatalogMapper.class);


//    @Mapping(target = "name", constant = "测试条目")
    ArticleCatalogDto toDto(ArticleCatalogDoc doc);

    ArticleCatalogDoc toDocument(ArticleCatalogDto doc);
}
