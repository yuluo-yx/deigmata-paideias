# crd-ref-docs

## crd 介绍

crd：CustomResourceDefinition（CRD），是Kubernetes（v1.7+）为提高可扩展性，让开发者去自定义资源的一种方式。

## elastic/crd-ref-docs 项目介绍

这是一个受https://github.com/ahmetb/gen-crd-api-reference-docs项目启发的全新实现。在尝试采用Kubernetes 上gen-crd-api-refernce-docs的 Elastic Cloud 生成文档时，我们遇到了一些缺点，例如缺乏对 Go 模块的支持、扫描速度慢以及渲染逻辑难以适应 Asciidoc（我们的首选文档标记语言） 。该项目试图通过重新实现类型发现逻辑并解耦渲染逻辑来解决这些问题，以便可以支持不同的标记格式。

## 测试 crd-ref-docs 功能

二进制包已经下载存放在当前路径中。（在 win 下编译）

```shell
# envoy-gateway: gen crds.

crd-ref-docs \
--source-path=api/v1alpha1 \
--config=tools/crd-ref-docs/config.yaml \
--templates-dir=tools/crd-ref-docs/templates \
--output-path=site/content/en/latest/api/extension_types.md \
--max-depth 10 \
--renderer=markdown

# test output:

crd-ref-docs \
--source-path=./example/api \
--config=./example/config.yaml \
--max-depth 10 \
--renderer=markdown \
--output-path=example/expected.md \
--templates-dir=templates/markdown

crd-ref-docs \
--source-path=test/api/ \
--config=test/config.yaml \
--max-depth 10 \
--renderer=asciidoctor \
--output-path=test/expected.asciidoc \
--templates-dir=./templates/asciidoctor

```
