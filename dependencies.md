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

| Dependency                     | License      |
| ------------------------------ | ------------ |
| github.com/DATA-DOG/go-sqlmock | [Unknown][9] |
| github.com/kinbiko/jsonassert  | [MIT][10]    |

## Extension Manager Java Client

### Compile Dependencies

| Dependency                      | License                                                                                                                                                                                             |
| ------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [swagger-annotations][11]       | [Apache License 2.0][12]                                                                                                                                                                            |
| [jersey-core-client][13]        | [EPL 2.0][14]; [GPL2 w/ CPE][15]; [EDL 1.0][16]; [BSD 2-Clause][17]; [Apache License, 2.0][12]; [Public Domain][18]; [Modified BSD][19]; [jQuery license][20]; [MIT license][21]; [W3C license][22] |
| [jersey-media-multipart][23]    | [EPL 2.0][14]; [GPL2 w/ CPE][15]; [EDL 1.0][16]; [BSD 2-Clause][17]; [Apache License, 2.0][12]; [Public Domain][18]; [Modified BSD][19]; [jQuery license][20]; [MIT license][21]; [W3C license][22] |
| [jersey-media-json-jackson][24] | [EPL 2.0][14]; [The GNU General Public License (GPL), Version 2, With Classpath Exception][15]; [Apache License, 2.0][12]                                                                           |
| [jersey-inject-hk2][25]         | [EPL 2.0][14]; [GPL2 w/ CPE][15]; [EDL 1.0][16]; [BSD 2-Clause][17]; [Apache License, 2.0][12]; [Public Domain][18]; [Modified BSD][19]; [jQuery license][20]; [MIT license][21]; [W3C license][22] |
| [Jackson-core][26]              | [The Apache Software License, Version 2.0][27]                                                                                                                                                      |
| [Jackson-annotations][28]       | [The Apache Software License, Version 2.0][27]                                                                                                                                                      |
| [jackson-databind][28]          | [The Apache Software License, Version 2.0][27]                                                                                                                                                      |
| [MiG Base64][29]                | [Prior BSD License][30]                                                                                                                                                                             |

### Test Dependencies

| Dependency                                | License                           |
| ----------------------------------------- | --------------------------------- |
| [JUnit Jupiter API][31]                   | [Eclipse Public License v2.0][32] |
| [EqualsVerifier | release normal jar][33] | [Apache License, Version 2.0][27] |

### Plugin Dependencies

| Dependency                                              | License                                        |
| ------------------------------------------------------- | ---------------------------------------------- |
| [Maven Flatten Plugin][34]                              | [Apache Software Licenese][27]                 |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][35] | [ASL2][36]                                     |
| [SonarQube Scanner for Maven][37]                       | [GNU LGPL 3][38]                               |
| [Apache Maven Compiler Plugin][39]                      | [Apache-2.0][27]                               |
| [Apache Maven Enforcer Plugin][40]                      | [Apache-2.0][27]                               |
| [Maven Surefire Plugin][41]                             | [Apache-2.0][27]                               |
| [Versions Maven Plugin][42]                             | [Apache License, Version 2.0][27]              |
| [duplicate-finder-maven-plugin Maven Mojo][43]          | [Apache License 2.0][12]                       |
| [Apache Maven Deploy Plugin][44]                        | [Apache-2.0][27]                               |
| [Apache Maven GPG Plugin][45]                           | [Apache License, Version 2.0][27]              |
| [Apache Maven Source Plugin][46]                        | [Apache License, Version 2.0][27]              |
| [Project keeper maven plugin][47]                       | [The MIT License][48]                          |
| [Exec Maven Plugin][49]                                 | [Apache License 2][27]                         |
| [swagger-codegen (maven-plugin)][50]                    | [Apache License 2.0][12]                       |
| [Build Helper Maven Plugin][51]                         | [The MIT License][52]                          |
| [Apache Maven Javadoc Plugin][53]                       | [Apache-2.0][27]                               |
| [Nexus Staging Maven Plugin][54]                        | [Eclipse Public License][55]                   |
| [JaCoCo :: Maven Plugin][56]                            | [Eclipse Public License 2.0][57]               |
| [error-code-crawler-maven-plugin][58]                   | [MIT License][59]                              |
| [Reproducible Build Maven Plugin][60]                   | [Apache 2.0][36]                               |
| [Maven Clean Plugin][61]                                | [The Apache Software License, Version 2.0][36] |
| [Maven Resources Plugin][62]                            | [The Apache Software License, Version 2.0][36] |
| [Maven JAR Plugin][63]                                  | [The Apache Software License, Version 2.0][36] |
| [Maven Install Plugin][64]                              | [The Apache Software License, Version 2.0][36] |
| [Maven Site Plugin 3][65]                               | [The Apache Software License, Version 2.0][36] |

## Extension Integration Tests Library

### Compile Dependencies

| Dependency                               | License                           |
| ---------------------------------------- | --------------------------------- |
| [Extension Manager Java Client][66]      | [MIT License][67]                 |
| [exasol-test-setup-abstraction-java][68] | [MIT License][69]                 |
| [Netty/Handler][70]                      | [Apache License, Version 2.0][71] |
| [Test Database Builder for Java][72]     | [MIT License][73]                 |
| [Matcher for SQL Result Sets][74]        | [MIT License][75]                 |
| [JUnit Jupiter API][31]                  | [Eclipse Public License v2.0][32] |

### Test Dependencies

| Dependency                  | License                           |
| --------------------------- | --------------------------------- |
| [JUnit Jupiter Params][31]  | [Eclipse Public License v2.0][32] |
| [mockito-junit-jupiter][76] | [The MIT License][77]             |
| [udf-debugging-java][78]    | [MIT License][79]                 |
| [SLF4J JDK14 Binding][80]   | [MIT License][21]                 |

### Plugin Dependencies

| Dependency                                              | License                                        |
| ------------------------------------------------------- | ---------------------------------------------- |
| [Maven Flatten Plugin][34]                              | [Apache Software Licenese][27]                 |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][35] | [ASL2][36]                                     |
| [SonarQube Scanner for Maven][37]                       | [GNU LGPL 3][38]                               |
| [Apache Maven Compiler Plugin][39]                      | [Apache-2.0][27]                               |
| [Apache Maven Enforcer Plugin][40]                      | [Apache-2.0][27]                               |
| [Maven Surefire Plugin][41]                             | [Apache-2.0][27]                               |
| [Versions Maven Plugin][42]                             | [Apache License, Version 2.0][27]              |
| [duplicate-finder-maven-plugin Maven Mojo][43]          | [Apache License 2.0][12]                       |
| [Apache Maven Deploy Plugin][44]                        | [Apache-2.0][27]                               |
| [Apache Maven GPG Plugin][45]                           | [Apache License, Version 2.0][27]              |
| [Apache Maven Source Plugin][46]                        | [Apache License, Version 2.0][27]              |
| [Apache Maven Javadoc Plugin][53]                       | [Apache-2.0][27]                               |
| [Nexus Staging Maven Plugin][54]                        | [Eclipse Public License][55]                   |
| [Maven Failsafe Plugin][81]                             | [Apache-2.0][27]                               |
| [JaCoCo :: Maven Plugin][56]                            | [Eclipse Public License 2.0][57]               |
| [error-code-crawler-maven-plugin][58]                   | [MIT License][59]                              |
| [Reproducible Build Maven Plugin][60]                   | [Apache 2.0][36]                               |
| [Project keeper maven plugin][47]                       | [The MIT License][48]                          |
| [Maven Clean Plugin][61]                                | [The Apache Software License, Version 2.0][36] |
| [Maven Resources Plugin][62]                            | [The Apache Software License, Version 2.0][36] |
| [Maven JAR Plugin][63]                                  | [The Apache Software License, Version 2.0][36] |
| [Maven Install Plugin][64]                              | [The Apache Software License, Version 2.0][36] |
| [Maven Site Plugin 3][65]                               | [The Apache Software License, Version 2.0][36] |

## Registry

### Compile Dependencies

| Dependency               | License          |
| ------------------------ | ---------------- |
| [aws-cdk-lib][82]        | [Apache-2.0][83] |
| [constructs][84]         | [Apache-2.0][85] |
| [source-map-support][86] | [MIT][87]        |

## Parametervalidator

### Compile Dependencies

| Dependency                                  | License |
| ------------------------------------------- | ------- |
| [@exasol/extension-parameter-validator][88] | MIT     |

[0]: https://github.com/Nightapes/go-rest/blob/v0.3.3/LICENSE
[1]: https://github.com/dop251/goja/blob/7749907a8a20/LICENSE
[2]: https://github.com/dop251/goja_nodejs/blob/804a84515562/LICENSE
[3]: https://github.com/exasol/exasol-driver-go/blob/v1.0.0/LICENSE
[4]: https://github.com/exasol/exasol-test-setup-abstraction-server/blob/go-client/v0.3.2/go-client/LICENSE
[5]: https://github.com/go-chi/chi/blob/v5.0.8/LICENSE
[6]: https://github.com/sirupsen/logrus/blob/v1.9.3/LICENSE
[7]: https://github.com/stretchr/testify/blob/v1.8.4/LICENSE
[8]: https://github.com/swaggo/http-swagger/blob/v1.3.4/LICENSE
[9]: https://github.com/DATA-DOG/go-sqlmock/blob/master/LICENSE
[10]: https://github.com/kinbiko/jsonassert/blob/HEAD/LICENSE
[11]: https://github.com/swagger-api/swagger-core/tree/master/modules/swagger-annotations
[12]: http://www.apache.org/licenses/LICENSE-2.0.html
[13]: https://projects.eclipse.org/projects/ee4j.jersey/jersey-client
[14]: http://www.eclipse.org/legal/epl-2.0
[15]: https://www.gnu.org/software/classpath/license.html
[16]: http://www.eclipse.org/org/documents/edl-v10.php
[17]: https://opensource.org/licenses/BSD-2-Clause
[18]: https://creativecommons.org/publicdomain/zero/1.0/
[19]: https://asm.ow2.io/license.html
[20]: https://jquery.org/license/
[21]: http://www.opensource.org/licenses/mit-license.php
[22]: https://www.w3.org/Consortium/Legal/copyright-documents-19990405
[23]: https://projects.eclipse.org/projects/ee4j.jersey/project/jersey-media-multipart
[24]: https://projects.eclipse.org/projects/ee4j.jersey/project/jersey-media-json-jackson
[25]: https://projects.eclipse.org/projects/ee4j.jersey/project/jersey-hk2
[26]: https://github.com/FasterXML/jackson-core
[27]: https://www.apache.org/licenses/LICENSE-2.0.txt
[28]: https://github.com/FasterXML/jackson
[29]: http://sourceforge.net/projects/migbase64/
[30]: http://en.wikipedia.org/wiki/BSD_licenses
[31]: https://junit.org/junit5/
[32]: https://www.eclipse.org/legal/epl-v20.html
[33]: https://www.jqno.nl/equalsverifier
[34]: https://www.mojohaus.org/flatten-maven-plugin/
[35]: https://sonatype.github.io/ossindex-maven/maven-plugin/
[36]: http://www.apache.org/licenses/LICENSE-2.0.txt
[37]: http://sonarsource.github.io/sonar-scanner-maven/
[38]: http://www.gnu.org/licenses/lgpl.txt
[39]: https://maven.apache.org/plugins/maven-compiler-plugin/
[40]: https://maven.apache.org/enforcer/maven-enforcer-plugin/
[41]: https://maven.apache.org/surefire/maven-surefire-plugin/
[42]: https://www.mojohaus.org/versions/versions-maven-plugin/
[43]: https://github.com/basepom/duplicate-finder-maven-plugin
[44]: https://maven.apache.org/plugins/maven-deploy-plugin/
[45]: https://maven.apache.org/plugins/maven-gpg-plugin/
[46]: https://maven.apache.org/plugins/maven-source-plugin/
[47]: https://github.com/exasol/project-keeper/
[48]: https://github.com/exasol/project-keeper/blob/main/LICENSE
[49]: https://www.mojohaus.org/exec-maven-plugin
[50]: https://github.com/swagger-api/swagger-codegen/tree/master/modules/swagger-codegen-maven-plugin
[51]: http://www.mojohaus.org/build-helper-maven-plugin/
[52]: https://opensource.org/licenses/mit-license.php
[53]: https://maven.apache.org/plugins/maven-javadoc-plugin/
[54]: http://www.sonatype.com/public-parent/nexus-maven-plugins/nexus-staging/nexus-staging-maven-plugin/
[55]: http://www.eclipse.org/legal/epl-v10.html
[56]: https://www.jacoco.org/jacoco/trunk/doc/maven.html
[57]: https://www.eclipse.org/legal/epl-2.0/
[58]: https://github.com/exasol/error-code-crawler-maven-plugin/
[59]: https://github.com/exasol/error-code-crawler-maven-plugin/blob/main/LICENSE
[60]: http://zlika.github.io/reproducible-build-maven-plugin
[61]: http://maven.apache.org/plugins/maven-clean-plugin/
[62]: http://maven.apache.org/plugins/maven-resources-plugin/
[63]: http://maven.apache.org/plugins/maven-jar-plugin/
[64]: http://maven.apache.org/plugins/maven-install-plugin/
[65]: http://maven.apache.org/plugins/maven-site-plugin/
[66]: https://github.com/exasol/extension-manager/
[67]: https://github.com/exasol/extension-manager/blob/main/LICENSE
[68]: https://github.com/exasol/exasol-test-setup-abstraction-java/
[69]: https://github.com/exasol/exasol-test-setup-abstraction-java/blob/main/LICENSE
[70]: https://netty.io/netty-handler/
[71]: https://www.apache.org/licenses/LICENSE-2.0
[72]: https://github.com/exasol/test-db-builder-java/
[73]: https://github.com/exasol/test-db-builder-java/blob/main/LICENSE
[74]: https://github.com/exasol/hamcrest-resultset-matcher/
[75]: https://github.com/exasol/hamcrest-resultset-matcher/blob/main/LICENSE
[76]: https://github.com/mockito/mockito
[77]: https://github.com/mockito/mockito/blob/main/LICENSE
[78]: https://github.com/exasol/udf-debugging-java/
[79]: https://github.com/exasol/udf-debugging-java/blob/main/LICENSE
[80]: http://www.slf4j.org
[81]: https://maven.apache.org/surefire/maven-failsafe-plugin/
[82]: https://registry.npmjs.org/aws-cdk-lib/-/aws-cdk-lib-2.72.1.tgz
[83]: https://github.com/aws/aws-cdk
[84]: https://registry.npmjs.org/constructs/-/constructs-10.1.300.tgz
[85]: https://github.com/aws/constructs
[86]: https://registry.npmjs.org/source-map-support/-/source-map-support-0.5.21.tgz
[87]: https://github.com/evanw/node-source-map-support
[88]: https://registry.npmjs.org/@exasol/extension-parameter-validator/-/extension-parameter-validator-0.2.0.tgz
