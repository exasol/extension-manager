<!-- @formatter:off -->
# Dependencies

## Extension-manager

### Compile Dependencies

| Dependency                                                       | License  |
| ---------------------------------------------------------------- | -------- |
| github.com/dop251/goja                                           | [MIT][0] |
| github.com/dop251/goja_nodejs                                    | [MIT][1] |
| github.com/exasol/exasol-driver-go                               | [MIT][2] |
| github.com/exasol/exasol-test-setup-abstraction-server/go-client | [MIT][3] |
| github.com/go-chi/chi/v5                                         | [MIT][4] |
| github.com/sirupsen/logrus                                       | [MIT][5] |
| github.com/stretchr/testify                                      | [MIT][6] |
| github.com/swaggo/http-swagger                                   | [MIT][7] |

### Test Dependencies

| Dependency                     | License      |
| ------------------------------ | ------------ |
| github.com/DATA-DOG/go-sqlmock | [Unknown][8] |
| github.com/Nightapes/go-rest   | [MIT][9]     |
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

| Dependency                                              | License                           |
| ------------------------------------------------------- | --------------------------------- |
| [Maven Flatten Plugin][34]                              | [Apache Software Licenese][27]    |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][35] | [ASL2][36]                        |
| [SonarQube Scanner for Maven][37]                       | [GNU LGPL 3][38]                  |
| [Apache Maven Compiler Plugin][39]                      | [Apache License, Version 2.0][27] |
| [Apache Maven Enforcer Plugin][40]                      | [Apache License, Version 2.0][27] |
| [Maven Surefire Plugin][41]                             | [Apache License, Version 2.0][27] |
| [Versions Maven Plugin][42]                             | [Apache License, Version 2.0][27] |
| [Apache Maven Deploy Plugin][43]                        | [Apache License, Version 2.0][27] |
| [Apache Maven GPG Plugin][44]                           | [Apache License, Version 2.0][27] |
| [Apache Maven Source Plugin][45]                        | [Apache License, Version 2.0][27] |
| [Project keeper maven plugin][46]                       | [The MIT License][47]             |
| [Exec Maven Plugin][48]                                 | [Apache License 2][27]            |
| [swagger-codegen (maven-plugin)][49]                    | [Apache License 2.0][12]          |
| [Build Helper Maven Plugin][50]                         | [The MIT License][51]             |
| [Apache Maven Javadoc Plugin][52]                       | [Apache License, Version 2.0][27] |
| [Nexus Staging Maven Plugin][53]                        | [Eclipse Public License][54]      |
| [JaCoCo :: Maven Plugin][55]                            | [Eclipse Public License 2.0][56]  |
| [error-code-crawler-maven-plugin][57]                   | [MIT License][58]                 |
| [Reproducible Build Maven Plugin][59]                   | [Apache 2.0][36]                  |
| [Apache Maven Clean Plugin][60]                         | [Apache License, Version 2.0][27] |
| [Apache Maven Resources Plugin][61]                     | [Apache License, Version 2.0][27] |
| [Apache Maven JAR Plugin][62]                           | [Apache License, Version 2.0][27] |
| [Apache Maven Install Plugin][63]                       | [Apache License, Version 2.0][27] |
| [Apache Maven Site Plugin][64]                          | [Apache License, Version 2.0][27] |

## Extension Integration Tests Library

### Compile Dependencies

| Dependency                                 | License                           |
| ------------------------------------------ | --------------------------------- |
| [Extension Manager Java Client][65]        | [MIT License][66]                 |
| [exasol-test-setup-abstraction-java][67]   | [MIT License][68]                 |
| [Test containers for Exasol on Docker][69] | [MIT License][70]                 |
| [Test Database Builder for Java][71]       | [MIT License][72]                 |
| [Matcher for SQL Result Sets][73]          | [MIT License][74]                 |
| [JUnit Jupiter API][31]                    | [Eclipse Public License v2.0][32] |

### Test Dependencies

| Dependency                  | License                           |
| --------------------------- | --------------------------------- |
| [JUnit Jupiter Params][31]  | [Eclipse Public License v2.0][32] |
| [mockito-junit-jupiter][75] | [The MIT License][76]             |
| [udf-debugging-java][77]    | [MIT License][78]                 |

### Plugin Dependencies

| Dependency                                              | License                           |
| ------------------------------------------------------- | --------------------------------- |
| [Maven Flatten Plugin][34]                              | [Apache Software Licenese][27]    |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][35] | [ASL2][36]                        |
| [SonarQube Scanner for Maven][37]                       | [GNU LGPL 3][38]                  |
| [Apache Maven Compiler Plugin][39]                      | [Apache License, Version 2.0][27] |
| [Apache Maven Enforcer Plugin][40]                      | [Apache License, Version 2.0][27] |
| [Maven Surefire Plugin][41]                             | [Apache License, Version 2.0][27] |
| [Versions Maven Plugin][42]                             | [Apache License, Version 2.0][27] |
| [Apache Maven Deploy Plugin][43]                        | [Apache License, Version 2.0][27] |
| [Apache Maven GPG Plugin][44]                           | [Apache License, Version 2.0][27] |
| [Apache Maven Source Plugin][45]                        | [Apache License, Version 2.0][27] |
| [Apache Maven Javadoc Plugin][52]                       | [Apache License, Version 2.0][27] |
| [Nexus Staging Maven Plugin][53]                        | [Eclipse Public License][54]      |
| [Maven Failsafe Plugin][79]                             | [Apache License, Version 2.0][27] |
| [JaCoCo :: Maven Plugin][55]                            | [Eclipse Public License 2.0][56]  |
| [error-code-crawler-maven-plugin][57]                   | [MIT License][58]                 |
| [Reproducible Build Maven Plugin][59]                   | [Apache 2.0][36]                  |
| [Project keeper maven plugin][46]                       | [The MIT License][47]             |
| [Apache Maven Clean Plugin][60]                         | [Apache License, Version 2.0][27] |
| [Apache Maven Resources Plugin][61]                     | [Apache License, Version 2.0][27] |
| [Apache Maven JAR Plugin][62]                           | [Apache License, Version 2.0][27] |
| [Apache Maven Install Plugin][63]                       | [Apache License, Version 2.0][27] |
| [Apache Maven Site Plugin][64]                          | [Apache License, Version 2.0][27] |

## Registry

### Compile Dependencies

| Dependency               | License          |
| ------------------------ | ---------------- |
| [aws-cdk-lib][80]        | [Apache-2.0][81] |
| [constructs][82]         | [Apache-2.0][83] |
| [source-map-support][84] | [MIT][85]        |

## Parametervalidator

### Compile Dependencies

| Dependency                                  | License |
| ------------------------------------------- | ------- |
| [@exasol/extension-parameter-validator][86] | MIT     |

[0]: https://github.com/dop251/goja/blob/5460598cfa32/LICENSE
[1]: https://github.com/dop251/goja_nodejs/blob/2229640ea097/LICENSE
[2]: https://github.com/exasol/exasol-driver-go/blob/v0.4.6/LICENSE
[3]: https://github.com/exasol/exasol-test-setup-abstraction-server/blob/go-client/v0.3.2/go-client/LICENSE
[4]: https://github.com/go-chi/chi/blob/v5.0.8/LICENSE
[5]: https://github.com/sirupsen/logrus/blob/v1.9.0/LICENSE
[6]: https://github.com/stretchr/testify/blob/v1.8.1/LICENSE
[7]: https://github.com/swaggo/http-swagger/blob/v1.3.3/LICENSE
[8]: https://github.com/DATA-DOG/go-sqlmock/blob/master/LICENSE
[9]: https://github.com/Nightapes/go-rest/blob/v0.3.1/LICENSE
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
[43]: https://maven.apache.org/plugins/maven-deploy-plugin/
[44]: https://maven.apache.org/plugins/maven-gpg-plugin/
[45]: https://maven.apache.org/plugins/maven-source-plugin/
[46]: https://github.com/exasol/project-keeper/
[47]: https://github.com/exasol/project-keeper/blob/main/LICENSE
[48]: https://www.mojohaus.org/exec-maven-plugin
[49]: https://github.com/swagger-api/swagger-codegen/tree/master/modules/swagger-codegen-maven-plugin
[50]: http://www.mojohaus.org/build-helper-maven-plugin/
[51]: https://opensource.org/licenses/mit-license.php
[52]: https://maven.apache.org/plugins/maven-javadoc-plugin/
[53]: http://www.sonatype.com/public-parent/nexus-maven-plugins/nexus-staging/nexus-staging-maven-plugin/
[54]: http://www.eclipse.org/legal/epl-v10.html
[55]: https://www.jacoco.org/jacoco/trunk/doc/maven.html
[56]: https://www.eclipse.org/legal/epl-2.0/
[57]: https://github.com/exasol/error-code-crawler-maven-plugin/
[58]: https://github.com/exasol/error-code-crawler-maven-plugin/blob/main/LICENSE
[59]: http://zlika.github.io/reproducible-build-maven-plugin
[60]: https://maven.apache.org/plugins/maven-clean-plugin/
[61]: https://maven.apache.org/plugins/maven-resources-plugin/
[62]: https://maven.apache.org/plugins/maven-jar-plugin/
[63]: https://maven.apache.org/plugins/maven-install-plugin/
[64]: https://maven.apache.org/plugins/maven-site-plugin/
[65]: https://github.com/exasol/extension-manager/
[66]: https://github.com/exasol/extension-manager/blob/main/LICENSE
[67]: https://github.com/exasol/exasol-test-setup-abstraction-java/
[68]: https://github.com/exasol/exasol-test-setup-abstraction-java/blob/main/LICENSE
[69]: https://github.com/exasol/exasol-testcontainers/
[70]: https://github.com/exasol/exasol-testcontainers/blob/main/LICENSE
[71]: https://github.com/exasol/test-db-builder-java/
[72]: https://github.com/exasol/test-db-builder-java/blob/main/LICENSE
[73]: https://github.com/exasol/hamcrest-resultset-matcher/
[74]: https://github.com/exasol/hamcrest-resultset-matcher/blob/main/LICENSE
[75]: https://github.com/mockito/mockito
[76]: https://github.com/mockito/mockito/blob/main/LICENSE
[77]: https://github.com/exasol/udf-debugging-java/
[78]: https://github.com/exasol/udf-debugging-java/blob/main/LICENSE
[79]: https://maven.apache.org/surefire/maven-failsafe-plugin/
[80]: https://registry.npmjs.org/aws-cdk-lib/-/aws-cdk-lib-2.62.0.tgz
[81]: https://github.com/aws/aws-cdk
[82]: https://registry.npmjs.org/constructs/-/constructs-10.1.231.tgz
[83]: https://github.com/aws/constructs
[84]: https://registry.npmjs.org/source-map-support/-/source-map-support-0.5.21.tgz
[85]: https://github.com/evanw/node-source-map-support
[86]: https://registry.npmjs.org/@exasol/extension-parameter-validator/-/extension-parameter-validator-0.2.0.tgz
