<!-- @formatter:off -->
# Dependencies

## Extension-manager

### Compile Dependencies

| Dependency                                                       | License  |
| ---------------------------------------------------------------- | -------- |
| github.com/Nightapes/go-rest                                     | [MIT][0] |
| github.com/dop251/goja                                           | [MIT][1] |
| github.com/dop251/goja_nodejs                                    | [MIT][2] |
| github.com/exasol/exasol-driver-go                               | [MIT][3] |
| github.com/exasol/exasol-test-setup-abstraction-server/go-client | [MIT][4] |
| github.com/go-chi/chi/v5                                         | [MIT][5] |
| github.com/sirupsen/logrus                                       | [MIT][6] |
| github.com/stretchr/testify                                      | [MIT][7] |
| github.com/swaggo/http-swagger                                   | [MIT][8] |

### Test Dependencies

| Dependency                     | License            |
| ------------------------------ | ------------------ |
| github.com/DATA-DOG/go-sqlmock | [Unknown][9]       |
| github.com/kinbiko/jsonassert  | [MIT][10]          |
| golang.org/x/mod               | [BSD-3-Clause][11] |

## Extension Manager Java Client

### Compile Dependencies

| Dependency                      | License                                                                                                                                                                                             |
| ------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [swagger-annotations][12]       | [Apache License 2.0][13]                                                                                                                                                                            |
| [jersey-core-client][14]        | [EPL 2.0][15]; [GPL2 w/ CPE][16]; [EDL 1.0][17]; [BSD 2-Clause][18]; [Apache License, 2.0][13]; [Public Domain][19]; [Modified BSD][20]; [jQuery license][21]; [MIT license][22]; [W3C license][23] |
| [jersey-media-multipart][24]    | [EPL 2.0][15]; [GPL2 w/ CPE][16]; [EDL 1.0][17]; [BSD 2-Clause][18]; [Apache License, 2.0][13]; [Public Domain][19]; [Modified BSD][20]; [jQuery license][21]; [MIT license][22]; [W3C license][23] |
| [jersey-media-json-jackson][25] | [EPL 2.0][15]; [The GNU General Public License (GPL), Version 2, With Classpath Exception][16]; [Apache License, 2.0][13]                                                                           |
| [jersey-inject-hk2][26]         | [EPL 2.0][15]; [GPL2 w/ CPE][16]; [EDL 1.0][17]; [BSD 2-Clause][18]; [Apache License, 2.0][13]; [Public Domain][19]; [Modified BSD][20]; [jQuery license][21]; [MIT license][22]; [W3C license][23] |
| [Jackson-core][27]              | [The Apache Software License, Version 2.0][28]                                                                                                                                                      |
| [Jackson-annotations][29]       | [The Apache Software License, Version 2.0][28]                                                                                                                                                      |
| [jackson-databind][29]          | [The Apache Software License, Version 2.0][28]                                                                                                                                                      |
| [MiG Base64][30]                | [Prior BSD License][31]                                                                                                                                                                             |

### Test Dependencies

| Dependency                                 | License                           |
| ------------------------------------------ | --------------------------------- |
| [JUnit Jupiter API][32]                    | [Eclipse Public License v2.0][33] |
| [JUnit Jupiter Params][32]                 | [Eclipse Public License v2.0][33] |
| [EqualsVerifier \| release normal jar][34] | [Apache License, Version 2.0][28] |

### Plugin Dependencies

| Dependency                                              | License                           |
| ------------------------------------------------------- | --------------------------------- |
| [SonarQube Scanner for Maven][35]                       | [GNU LGPL 3][36]                  |
| [Apache Maven Toolchains Plugin][37]                    | [Apache License, Version 2.0][28] |
| [Apache Maven Compiler Plugin][38]                      | [Apache-2.0][28]                  |
| [Apache Maven Enforcer Plugin][39]                      | [Apache-2.0][28]                  |
| [Maven Flatten Plugin][40]                              | [Apache Software Licenese][28]    |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][41] | [ASL2][42]                        |
| [Maven Surefire Plugin][43]                             | [Apache-2.0][28]                  |
| [Versions Maven Plugin][44]                             | [Apache License, Version 2.0][28] |
| [duplicate-finder-maven-plugin Maven Mojo][45]          | [Apache License 2.0][13]          |
| [Apache Maven Deploy Plugin][46]                        | [Apache-2.0][28]                  |
| [Apache Maven GPG Plugin][47]                           | [Apache-2.0][28]                  |
| [Apache Maven Source Plugin][48]                        | [Apache License, Version 2.0][28] |
| [Exec Maven Plugin][49]                                 | [Apache License 2][28]            |
| [swagger-codegen (maven-plugin)][50]                    | [Apache License 2.0][13]          |
| [Build Helper Maven Plugin][51]                         | [The MIT License][52]             |
| [Apache Maven Javadoc Plugin][53]                       | [Apache-2.0][28]                  |
| [Nexus Staging Maven Plugin][54]                        | [Eclipse Public License][55]      |
| [JaCoCo :: Maven Plugin][56]                            | [Eclipse Public License 2.0][57]  |
| [error-code-crawler-maven-plugin][58]                   | [MIT License][59]                 |
| [Reproducible Build Maven Plugin][60]                   | [Apache 2.0][42]                  |

## Extension Integration Tests Library

### Compile Dependencies

| Dependency                               | License                           |
| ---------------------------------------- | --------------------------------- |
| [Extension Manager Java Client][61]      | [MIT License][62]                 |
| [exasol-test-setup-abstraction-java][63] | [MIT License][64]                 |
| [Test Database Builder for Java][65]     | [MIT License][66]                 |
| [Matcher for SQL Result Sets][67]        | [MIT License][68]                 |
| [JUnit Jupiter API][32]                  | [Eclipse Public License v2.0][33] |

### Test Dependencies

| Dependency                         | License                           |
| ---------------------------------- | --------------------------------- |
| [JUnit Jupiter Params][32]         | [Eclipse Public License v2.0][33] |
| [mockito-junit-jupiter][69]        | [MIT][70]                         |
| [udf-debugging-java][71]           | [MIT License][72]                 |
| [Maven Project Version Getter][73] | [MIT License][74]                 |
| [SLF4J JDK14 Provider][75]         | [MIT License][22]                 |

### Plugin Dependencies

| Dependency                                              | License                           |
| ------------------------------------------------------- | --------------------------------- |
| [SonarQube Scanner for Maven][35]                       | [GNU LGPL 3][36]                  |
| [Apache Maven Toolchains Plugin][37]                    | [Apache License, Version 2.0][28] |
| [Apache Maven Compiler Plugin][38]                      | [Apache-2.0][28]                  |
| [Apache Maven Enforcer Plugin][39]                      | [Apache-2.0][28]                  |
| [Maven Flatten Plugin][40]                              | [Apache Software Licenese][28]    |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][41] | [ASL2][42]                        |
| [Maven Surefire Plugin][43]                             | [Apache-2.0][28]                  |
| [Versions Maven Plugin][44]                             | [Apache License, Version 2.0][28] |
| [duplicate-finder-maven-plugin Maven Mojo][45]          | [Apache License 2.0][13]          |
| [Apache Maven Deploy Plugin][46]                        | [Apache-2.0][28]                  |
| [Apache Maven GPG Plugin][47]                           | [Apache-2.0][28]                  |
| [Apache Maven Source Plugin][48]                        | [Apache License, Version 2.0][28] |
| [Apache Maven Javadoc Plugin][53]                       | [Apache-2.0][28]                  |
| [Nexus Staging Maven Plugin][54]                        | [Eclipse Public License][55]      |
| [Maven Failsafe Plugin][76]                             | [Apache-2.0][28]                  |
| [JaCoCo :: Maven Plugin][56]                            | [Eclipse Public License 2.0][57]  |
| [error-code-crawler-maven-plugin][58]                   | [MIT License][59]                 |
| [Reproducible Build Maven Plugin][60]                   | [Apache 2.0][42]                  |
| [Apache Maven JAR Plugin][77]                           | [Apache License, Version 2.0][28] |

## Registry

### Compile Dependencies

| Dependency               | License          |
| ------------------------ | ---------------- |
| [aws-cdk-lib][78]        | [Apache-2.0][79] |
| [constructs][80]         | [Apache-2.0][81] |
| [source-map-support][82] | [MIT][83]        |

## Registry-upload

### Compile Dependencies

| Dependency                           | License          |
| ------------------------------------ | ---------------- |
| [@aws-sdk/client-cloudformation][84] | [Apache-2.0][85] |
| [@aws-sdk/client-cloudfront][86]     | [Apache-2.0][85] |
| [@aws-sdk/client-s3][87]             | [Apache-2.0][85] |
| [follow-redirects][88]               | [MIT][89]        |
| [octokit][90]                        | [MIT][91]        |

## Parametervalidator

### Compile Dependencies

| Dependency                                  | License |
| ------------------------------------------- | ------- |
| [@exasol/extension-parameter-validator][92] | MIT     |

[0]: https://github.com/Nightapes/go-rest/blob/v0.3.3/LICENSE
[1]: https://github.com/dop251/goja/blob/e401ed450204/LICENSE
[2]: https://github.com/dop251/goja_nodejs/blob/27eeffc9c235/LICENSE
[3]: https://github.com/exasol/exasol-driver-go/blob/v1.0.6/LICENSE
[4]: https://github.com/exasol/exasol-test-setup-abstraction-server/blob/go-client/v0.3.6/go-client/LICENSE
[5]: https://github.com/go-chi/chi/blob/v5.0.12/LICENSE
[6]: https://github.com/sirupsen/logrus/blob/v1.9.3/LICENSE
[7]: https://github.com/stretchr/testify/blob/v1.9.0/LICENSE
[8]: https://github.com/swaggo/http-swagger/blob/v1.3.4/LICENSE
[9]: https://github.com/DATA-DOG/go-sqlmock/blob/master/LICENSE
[10]: https://github.com/kinbiko/jsonassert/blob/HEAD/LICENSE
[11]: https://cs.opensource.google/go/x/mod/+/v0.17.0:LICENSE
[12]: https://github.com/swagger-api/swagger-core/tree/master/modules/swagger-annotations
[13]: http://www.apache.org/licenses/LICENSE-2.0.html
[14]: https://projects.eclipse.org/projects/ee4j.jersey/jersey-client
[15]: http://www.eclipse.org/legal/epl-2.0
[16]: https://www.gnu.org/software/classpath/license.html
[17]: http://www.eclipse.org/org/documents/edl-v10.php
[18]: https://opensource.org/licenses/BSD-2-Clause
[19]: https://creativecommons.org/publicdomain/zero/1.0/
[20]: https://asm.ow2.io/license.html
[21]: https://jquery.org/license/
[22]: http://www.opensource.org/licenses/mit-license.php
[23]: https://www.w3.org/Consortium/Legal/copyright-documents-19990405
[24]: https://projects.eclipse.org/projects/ee4j.jersey/project/jersey-media-multipart
[25]: https://projects.eclipse.org/projects/ee4j.jersey/project/jersey-media-json-jackson
[26]: https://projects.eclipse.org/projects/ee4j.jersey/project/jersey-hk2
[27]: https://github.com/FasterXML/jackson-core
[28]: https://www.apache.org/licenses/LICENSE-2.0.txt
[29]: https://github.com/FasterXML/jackson
[30]: http://sourceforge.net/projects/migbase64/
[31]: http://en.wikipedia.org/wiki/BSD_licenses
[32]: https://junit.org/junit5/
[33]: https://www.eclipse.org/legal/epl-v20.html
[34]: https://www.jqno.nl/equalsverifier
[35]: http://sonarsource.github.io/sonar-scanner-maven/
[36]: http://www.gnu.org/licenses/lgpl.txt
[37]: https://maven.apache.org/plugins/maven-toolchains-plugin/
[38]: https://maven.apache.org/plugins/maven-compiler-plugin/
[39]: https://maven.apache.org/enforcer/maven-enforcer-plugin/
[40]: https://www.mojohaus.org/flatten-maven-plugin/
[41]: https://sonatype.github.io/ossindex-maven/maven-plugin/
[42]: http://www.apache.org/licenses/LICENSE-2.0.txt
[43]: https://maven.apache.org/surefire/maven-surefire-plugin/
[44]: https://www.mojohaus.org/versions/versions-maven-plugin/
[45]: https://basepom.github.io/duplicate-finder-maven-plugin
[46]: https://maven.apache.org/plugins/maven-deploy-plugin/
[47]: https://maven.apache.org/plugins/maven-gpg-plugin/
[48]: https://maven.apache.org/plugins/maven-source-plugin/
[49]: https://www.mojohaus.org/exec-maven-plugin
[50]: https://github.com/swagger-api/swagger-codegen/tree/master/modules/swagger-codegen-maven-plugin
[51]: https://www.mojohaus.org/build-helper-maven-plugin/
[52]: https://spdx.org/licenses/MIT.txt
[53]: https://maven.apache.org/plugins/maven-javadoc-plugin/
[54]: http://www.sonatype.com/public-parent/nexus-maven-plugins/nexus-staging/nexus-staging-maven-plugin/
[55]: http://www.eclipse.org/legal/epl-v10.html
[56]: https://www.jacoco.org/jacoco/trunk/doc/maven.html
[57]: https://www.eclipse.org/legal/epl-2.0/
[58]: https://github.com/exasol/error-code-crawler-maven-plugin/
[59]: https://github.com/exasol/error-code-crawler-maven-plugin/blob/main/LICENSE
[60]: http://zlika.github.io/reproducible-build-maven-plugin
[61]: https://github.com/exasol/extension-manager/
[62]: https://github.com/exasol/extension-manager/blob/main/LICENSE
[63]: https://github.com/exasol/exasol-test-setup-abstraction-java/
[64]: https://github.com/exasol/exasol-test-setup-abstraction-java/blob/main/LICENSE
[65]: https://github.com/exasol/test-db-builder-java/
[66]: https://github.com/exasol/test-db-builder-java/blob/main/LICENSE
[67]: https://github.com/exasol/hamcrest-resultset-matcher/
[68]: https://github.com/exasol/hamcrest-resultset-matcher/blob/main/LICENSE
[69]: https://github.com/mockito/mockito
[70]: https://opensource.org/licenses/MIT
[71]: https://github.com/exasol/udf-debugging-java/
[72]: https://github.com/exasol/udf-debugging-java/blob/main/LICENSE
[73]: https://github.com/exasol/maven-project-version-getter/
[74]: https://github.com/exasol/maven-project-version-getter/blob/main/LICENSE
[75]: http://www.slf4j.org
[76]: https://maven.apache.org/surefire/maven-failsafe-plugin/
[77]: https://maven.apache.org/plugins/maven-jar-plugin/
[78]: https://registry.npmjs.org/aws-cdk-lib/-/aws-cdk-lib-2.133.0.tgz
[79]: https://github.com/aws/aws-cdk
[80]: https://registry.npmjs.org/constructs/-/constructs-10.3.0.tgz
[81]: https://github.com/aws/constructs
[82]: https://registry.npmjs.org/source-map-support/-/source-map-support-0.5.21.tgz
[83]: https://github.com/evanw/node-source-map-support
[84]: https://registry.npmjs.org/@aws-sdk/client-cloudformation/-/client-cloudformation-3.537.0.tgz
[85]: https://github.com/aws/aws-sdk-js-v3
[86]: https://registry.npmjs.org/@aws-sdk/client-cloudfront/-/client-cloudfront-3.535.0.tgz
[87]: https://registry.npmjs.org/@aws-sdk/client-s3/-/client-s3-3.537.0.tgz
[88]: https://registry.npmjs.org/follow-redirects/-/follow-redirects-1.15.6.tgz
[89]: https://github.com/follow-redirects/follow-redirects
[90]: https://registry.npmjs.org/octokit/-/octokit-3.1.2.tgz
[91]: https://github.com/octokit/octokit.js
[92]: https://registry.npmjs.org/@exasol/extension-parameter-validator/-/extension-parameter-validator-0.3.0.tgz
