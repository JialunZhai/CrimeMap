workspace(name = "edu_nyu_cs_crimemap")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Go

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "6b65cb7917b4d1709f9410ffe00ecf3e160edf674b78c54a894471320862184f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.39.0/rules_go-v0.39.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.39.0/rules_go-v0.39.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "727f3e4edd96ea20c29e8c2ca9e8d2af724d8c7778e7923a854b2c80952bc405",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.30.0/bazel-gazelle-v0.30.0.tar.gz",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

go_repository(
    name = "com_github_trinodb_trino_go_client",
    importpath = "github.com/trinodb/trino-go-client",
    sum = "h1:GjSG/60MdmaZHWmOsUAbpwjElCcgfoem6KvFWeJ0Hss=",
    version = "v0.310.0",
)

go_repository(
    name = "in_gopkg_jcmturner_gokrb5_v6",
    importpath = "gopkg.in/jcmturner/gokrb5.v6",
    sum = "h1:n0KFjpbuM5pFMN38/Ay+Br3l91netGSVqHPHEXeWUqk=",
    version = "v6.1.1",
)

go_repository(
    name = "com_github_jcmturner_gofork",
    importpath = "github.com/jcmturner/gofork",
    sum = "h1:QH0l3hzAU1tfT3rZCnW5zXl+orbkNMMRGJfdJjHVETg=",
    version = "v1.7.6",
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    sum = "h1:AvwMYaRytfdeVt3u6mLaxYtErKYjxA2OXjJ1HHq6t3A=",
    version = "v0.7.0",
)

go_repository(
    name = "in_gopkg_jcmturner_dnsutils_v1",
    importpath = "gopkg.in/jcmturner/dnsutils.v1",
    sum = "h1:cIuC1OLRGZrld+16ZJvvZxVJeKPsvd5eUIvxfoN5hSM=",
    version = "v1.0.1",
)

go_repository(
    name = "in_gopkg_jcmturner_aescts_v1",
    importpath = "gopkg.in/jcmturner/aescts.v1",
    sum = "h1:cVVZBK2b1zY26haWB4vbBiZrfFQnfbTVrE3xZq6hrEw=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_hashicorp_go_uuid",
    importpath = "github.com/hashicorp/go-uuid",
    sum = "h1:2gKiV6YVmrJ1i2CKKa9obLvRieoRGviZFL26PcT/Co8=",
    version = "v1.0.3",
)

go_rules_dependencies()

go_register_toolchains(version = "1.20.2")

gazelle_dependencies()

# protobuf

http_archive(
    name = "com_google_protobuf",
    sha256 = "ba0650be1b169d24908eeddbe6107f011d8df0da5b1a5a4449a913b10e578faf",
    strip_prefix = "protobuf-3.19.4",
    urls = ["https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protobuf-all-3.19.4.tar.gz"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

# webapp

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "94070eff79305be05b7699207fbac5d2608054dd53e6109f7d00d923919ff45a",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/5.8.2/rules_nodejs-5.8.2.tar.gz"],
)

load("@build_bazel_rules_nodejs//:repositories.bzl", "build_bazel_rules_nodejs_dependencies")

build_bazel_rules_nodejs_dependencies()

load("@build_bazel_rules_nodejs//:index.bzl", "yarn_install")

yarn_install(
    name = "npm_deps",
    # Cannot opt-in yet because protractor rule looks for
    # @npm//:node_modules/protractor/bin/protractor
    exports_directories_only = False,
    package_json = "//:package.json",
    yarn_lock = "//:yarn.lock",
)

load("@npm_deps//@bazel/protractor:package.bzl", "npm_bazel_protractor_dependencies")

npm_bazel_protractor_dependencies()

# Setup the rules_webtesting toolchain
load("@io_bazel_rules_webtesting//web:repositories.bzl", "web_test_repositories")

web_test_repositories()

load("@io_bazel_rules_webtesting//web/versioned:browsers-0.3.3.bzl", "browser_repositories")

browser_repositories(
    chromium = True,
    firefox = True,
)
