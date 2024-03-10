package bookstore.queryapi.mapper;

import bookstore.queryapi.document.ArticleCatalogDoc;
import bookstore.queryapi.dto.ArticleCatalogDto;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.ReportingPolicy;
import org.mapstruct.factory.Mappers;

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface ArticleCatalogMapper {
    ArticleCatalogMapper INSTANCE = Mappers.getMapper(ArticleCatalogMapper.class);


    //    @Mapping(target = "name", constant = "测试条目")
    @Mapping(target = "name")
    ArticleCatalogDto toDto(ArticleCatalogDoc doc);

    ArticleCatalogDoc toDocument(ArticleCatalogDto doc);
}
