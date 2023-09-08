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

| Dependency                                              | License                                        |
| ------------------------------------------------------- | ---------------------------------------------- |
| [Maven Flatten Plugin][35]                              | [Apache Software Licenese][28]                 |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][36] | [ASL2][37]                                     |
| [SonarQube Scanner for Maven][38]                       | [GNU LGPL 3][39]                               |
| [Apache Maven Compiler Plugin][40]                      | [Apache-2.0][28]                               |
| [Apache Maven Enforcer Plugin][41]                      | [Apache-2.0][28]                               |
| [Maven Surefire Plugin][42]                             | [Apache-2.0][28]                               |
| [Versions Maven Plugin][43]                             | [Apache License, Version 2.0][28]              |
| [duplicate-finder-maven-plugin Maven Mojo][44]          | [Apache License 2.0][13]                       |
| [Apache Maven Deploy Plugin][45]                        | [Apache-2.0][28]                               |
| [Apache Maven GPG Plugin][46]                           | [Apache-2.0][28]                               |
| [Apache Maven Source Plugin][47]                        | [Apache License, Version 2.0][28]              |
| [Project keeper maven plugin][48]                       | [The MIT License][49]                          |
| [Exec Maven Plugin][50]                                 | [Apache License 2][28]                         |
| [swagger-codegen (maven-plugin)][51]                    | [Apache License 2.0][13]                       |
| [Build Helper Maven Plugin][52]                         | [The MIT License][53]                          |
| [Apache Maven Javadoc Plugin][54]                       | [Apache-2.0][28]                               |
| [Nexus Staging Maven Plugin][55]                        | [Eclipse Public License][56]                   |
| [JaCoCo :: Maven Plugin][57]                            | [Eclipse Public License 2.0][58]               |
| [error-code-crawler-maven-plugin][59]                   | [MIT License][60]                              |
| [Reproducible Build Maven Plugin][61]                   | [Apache 2.0][37]                               |
| [Maven Clean Plugin][62]                                | [The Apache Software License, Version 2.0][37] |
| [Maven Resources Plugin][63]                            | [The Apache Software License, Version 2.0][37] |
| [Maven JAR Plugin][64]                                  | [The Apache Software License, Version 2.0][37] |
| [Maven Install Plugin][65]                              | [The Apache Software License, Version 2.0][37] |
| [Maven Site Plugin 3][66]                               | [The Apache Software License, Version 2.0][37] |

## Extension Integration Tests Library

### Compile Dependencies

| Dependency                               | License                           |
| ---------------------------------------- | --------------------------------- |
| [Extension Manager Java Client][67]      | [MIT License][68]                 |
| [exasol-test-setup-abstraction-java][69] | [MIT License][70]                 |
| [Test Database Builder for Java][71]     | [MIT License][72]                 |
| [Matcher for SQL Result Sets][73]        | [MIT License][74]                 |
| [JUnit Jupiter API][32]                  | [Eclipse Public License v2.0][33] |

### Test Dependencies

| Dependency                         | License                           |
| ---------------------------------- | --------------------------------- |
| [JUnit Jupiter Params][32]         | [Eclipse Public License v2.0][33] |
| [mockito-junit-jupiter][75]        | [The MIT License][76]             |
| [udf-debugging-java][77]           | [MIT License][78]                 |
| [Maven Project Version Getter][79] | [MIT License][80]                 |
| [SLF4J JDK14 Provider][81]         | [MIT License][22]                 |

### Plugin Dependencies

| Dependency                                              | License                                        |
| ------------------------------------------------------- | ---------------------------------------------- |
| [Maven Flatten Plugin][35]                              | [Apache Software Licenese][28]                 |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][36] | [ASL2][37]                                     |
| [SonarQube Scanner for Maven][38]                       | [GNU LGPL 3][39]                               |
| [Apache Maven Compiler Plugin][40]                      | [Apache-2.0][28]                               |
| [Apache Maven Enforcer Plugin][41]                      | [Apache-2.0][28]                               |
| [Maven Surefire Plugin][42]                             | [Apache-2.0][28]                               |
| [Versions Maven Plugin][43]                             | [Apache License, Version 2.0][28]              |
| [duplicate-finder-maven-plugin Maven Mojo][44]          | [Apache License 2.0][13]                       |
| [Apache Maven Deploy Plugin][45]                        | [Apache-2.0][28]                               |
| [Apache Maven GPG Plugin][46]                           | [Apache-2.0][28]                               |
| [Apache Maven Source Plugin][47]                        | [Apache License, Version 2.0][28]              |
| [Apache Maven Javadoc Plugin][54]                       | [Apache-2.0][28]                               |
| [Nexus Staging Maven Plugin][55]                        | [Eclipse Public License][56]                   |
| [Maven Failsafe Plugin][82]                             | [Apache-2.0][28]                               |
| [JaCoCo :: Maven Plugin][57]                            | [Eclipse Public License 2.0][58]               |
| [error-code-crawler-maven-plugin][59]                   | [MIT License][60]                              |
| [Reproducible Build Maven Plugin][61]                   | [Apache 2.0][37]                               |
| [Project keeper maven plugin][48]                       | [The MIT License][49]                          |
| [Apache Maven JAR Plugin][83]                           | [Apache License, Version 2.0][28]              |
| [Maven Clean Plugin][62]                                | [The Apache Software License, Version 2.0][37] |
| [Maven Resources Plugin][63]                            | [The Apache Software License, Version 2.0][37] |
| [Maven Install Plugin][65]                              | [The Apache Software License, Version 2.0][37] |
| [Maven Site Plugin 3][66]                               | [The Apache Software License, Version 2.0][37] |

## Registry

### Compile Dependencies

| Dependency               | License          |
| ------------------------ | ---------------- |
| [aws-cdk-lib][84]        | [Apache-2.0][85] |
| [constructs][86]         | [Apache-2.0][87] |
| [source-map-support][88] | [MIT][89]        |

## Parametervalidator

### Compile Dependencies

| Dependency                                  | License |
| ------------------------------------------- | ------- |
| [@exasol/extension-parameter-validator][90] | MIT     |

[0]: https://github.com/Nightapes/go-rest/blob/v0.3.3/LICENSE
[1]: https://github.com/dop251/goja/blob/636fdf960de1/LICENSE
[2]: https://github.com/dop251/goja_nodejs/blob/804a84515562/LICENSE
[3]: https://github.com/exasol/exasol-driver-go/blob/v1.0.0/LICENSE
[4]: https://github.com/exasol/exasol-test-setup-abstraction-server/blob/go-client/v0.3.3/go-client/LICENSE
[5]: https://github.com/go-chi/chi/blob/v5.0.10/LICENSE
[6]: https://github.com/sirupsen/logrus/blob/v1.9.3/LICENSE
[7]: https://github.com/stretchr/testify/blob/v1.8.4/LICENSE
[8]: https://github.com/swaggo/http-swagger/blob/v1.3.4/LICENSE
[9]: https://github.com/DATA-DOG/go-sqlmock/blob/master/LICENSE
[10]: https://github.com/kinbiko/jsonassert/blob/HEAD/LICENSE
[11]: https://cs.opensource.google/go/x/mod/+/v0.12.0:LICENSE
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
[35]: https://www.mojohaus.org/flatten-maven-plugin/
[36]: https://sonatype.github.io/ossindex-maven/maven-plugin/
[37]: http://www.apache.org/licenses/LICENSE-2.0.txt
[38]: http://sonarsource.github.io/sonar-scanner-maven/
[39]: http://www.gnu.org/licenses/lgpl.txt
[40]: https://maven.apache.org/plugins/maven-compiler-plugin/
[41]: https://maven.apache.org/enforcer/maven-enforcer-plugin/
[42]: https://maven.apache.org/surefire/maven-surefire-plugin/
[43]: https://www.mojohaus.org/versions/versions-maven-plugin/
[44]: https://basepom.github.io/duplicate-finder-maven-plugin
[45]: https://maven.apache.org/plugins/maven-deploy-plugin/
[46]: https://maven.apache.org/plugins/maven-gpg-plugin/
[47]: https://maven.apache.org/plugins/maven-source-plugin/
[48]: https://github.com/exasol/project-keeper/
[49]: https://github.com/exasol/project-keeper/blob/main/LICENSE
[50]: https://www.mojohaus.org/exec-maven-plugin
[51]: https://github.com/swagger-api/swagger-codegen/tree/master/modules/swagger-codegen-maven-plugin
[52]: https://www.mojohaus.org/build-helper-maven-plugin/
[53]: https://spdx.org/licenses/MIT.txt
[54]: https://maven.apache.org/plugins/maven-javadoc-plugin/
[55]: http://www.sonatype.com/public-parent/nexus-maven-plugins/nexus-staging/nexus-staging-maven-plugin/
[56]: http://www.eclipse.org/legal/epl-v10.html
[57]: https://www.jacoco.org/jacoco/trunk/doc/maven.html
[58]: https://www.eclipse.org/legal/epl-2.0/
[59]: https://github.com/exasol/error-code-crawler-maven-plugin/
[60]: https://github.com/exasol/error-code-crawler-maven-plugin/blob/main/LICENSE
[61]: http://zlika.github.io/reproducible-build-maven-plugin
[62]: http://maven.apache.org/plugins/maven-clean-plugin/
[63]: http://maven.apache.org/plugins/maven-resources-plugin/
[64]: http://maven.apache.org/plugins/maven-jar-plugin/
[65]: http://maven.apache.org/plugins/maven-install-plugin/
[66]: http://maven.apache.org/plugins/maven-site-plugin/
[67]: https://github.com/exasol/extension-manager/
[68]: https://github.com/exasol/extension-manager/blob/main/LICENSE
[69]: https://github.com/exasol/exasol-test-setup-abstraction-java/
[70]: https://github.com/exasol/exasol-test-setup-abstraction-java/blob/main/LICENSE
[71]: https://github.com/exasol/test-db-builder-java/
[72]: https://github.com/exasol/test-db-builder-java/blob/main/LICENSE
[73]: https://github.com/exasol/hamcrest-resultset-matcher/
[74]: https://github.com/exasol/hamcrest-resultset-matcher/blob/main/LICENSE
[75]: https://github.com/mockito/mockito
[76]: https://github.com/mockito/mockito/blob/main/LICENSE
[77]: https://github.com/exasol/udf-debugging-java/
[78]: https://github.com/exasol/udf-debugging-java/blob/main/LICENSE
[79]: https://github.com/exasol/maven-project-version-getter/
[80]: https://github.com/exasol/maven-project-version-getter/blob/main/LICENSE
[81]: http://www.slf4j.org
[82]: https://maven.apache.org/surefire/maven-failsafe-plugin/
[83]: https://maven.apache.org/plugins/maven-jar-plugin/
[84]: https://registry.npmjs.org/aws-cdk-lib/-/aws-cdk-lib-2.94.0.tgz
[85]: https://github.com/aws/aws-cdk
[86]: https://registry.npmjs.org/constructs/-/constructs-10.2.70.tgz
[87]: https://github.com/aws/constructs
[88]: https://registry.npmjs.org/source-map-support/-/source-map-support-0.5.21.tgz
[89]: https://github.com/evanw/node-source-map-support
[90]: https://registry.npmjs.org/@exasol/extension-parameter-validator/-/extension-parameter-validator-0.3.0.tgz
