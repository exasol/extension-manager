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

### Plugin Dependencies

| Dependency                                              | License                                        |
| ------------------------------------------------------- | ---------------------------------------------- |
| [SonarQube Scanner for Maven][31]                       | [GNU LGPL 3][32]                               |
| [Apache Maven Compiler Plugin][33]                      | [Apache License, Version 2.0][34]              |
| [Apache Maven Enforcer Plugin][35]                      | [Apache License, Version 2.0][34]              |
| [Maven Flatten Plugin][36]                              | [Apache Software Licenese][27]                 |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][37] | [ASL2][27]                                     |
| [Maven Surefire Plugin][38]                             | [Apache License, Version 2.0][34]              |
| [Versions Maven Plugin][39]                             | [Apache License, Version 2.0][34]              |
| [Apache Maven Deploy Plugin][40]                        | [Apache License, Version 2.0][34]              |
| [Apache Maven GPG Plugin][41]                           | [Apache License, Version 2.0][34]              |
| [Apache Maven Source Plugin][42]                        | [Apache License, Version 2.0][34]              |
| [Project keeper maven plugin][43]                       | [The MIT License][44]                          |
| [Exec Maven Plugin][45]                                 | [Apache License 2][27]                         |
| [swagger-codegen (maven-plugin)][46]                    | [Apache License 2.0][12]                       |
| [Build Helper Maven Plugin][47]                         | [The MIT License][48]                          |
| [Apache Maven Javadoc Plugin][49]                       | [Apache License, Version 2.0][34]              |
| [Nexus Staging Maven Plugin][50]                        | [Eclipse Public License][51]                   |
| [JaCoCo :: Maven Plugin][52]                            | [Eclipse Public License 2.0][53]               |
| [error-code-crawler-maven-plugin][54]                   | [MIT License][55]                              |
| [Reproducible Build Maven Plugin][56]                   | [Apache 2.0][27]                               |
| [Maven Clean Plugin][57]                                | [The Apache Software License, Version 2.0][27] |
| [Maven Resources Plugin][58]                            | [The Apache Software License, Version 2.0][27] |
| [Maven JAR Plugin][59]                                  | [The Apache Software License, Version 2.0][27] |
| [Maven Install Plugin][60]                              | [The Apache Software License, Version 2.0][27] |
| [Maven Site Plugin 3][61]                               | [The Apache Software License, Version 2.0][27] |

## Extension Integration Tests Library

### Compile Dependencies

| Dependency                               | License                           |
| ---------------------------------------- | --------------------------------- |
| [Extension Manager Java Client][62]      | [MIT License][63]                 |
| [exasol-test-setup-abstraction-java][64] | [MIT License][65]                 |
| [Test Database Builder for Java][66]     | [MIT License][67]                 |
| [Matcher for SQL Result Sets][68]        | [MIT License][69]                 |
| [JUnit Jupiter API][70]                  | [Eclipse Public License v2.0][71] |

### Test Dependencies

| Dependency                 | License                           |
| -------------------------- | --------------------------------- |
| [JUnit Jupiter Params][70] | [Eclipse Public License v2.0][71] |
| [mockito-core][72]         | [The MIT License][73]             |
| [udf-debugging-java][74]   | [MIT][75]                         |

### Plugin Dependencies

| Dependency                                              | License                                        |
| ------------------------------------------------------- | ---------------------------------------------- |
| [SonarQube Scanner for Maven][31]                       | [GNU LGPL 3][32]                               |
| [Apache Maven Compiler Plugin][33]                      | [Apache License, Version 2.0][34]              |
| [Apache Maven Enforcer Plugin][35]                      | [Apache License, Version 2.0][34]              |
| [Maven Flatten Plugin][36]                              | [Apache Software Licenese][27]                 |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][37] | [ASL2][27]                                     |
| [Maven Surefire Plugin][38]                             | [Apache License, Version 2.0][34]              |
| [Versions Maven Plugin][39]                             | [Apache License, Version 2.0][34]              |
| [Apache Maven Deploy Plugin][40]                        | [Apache License, Version 2.0][34]              |
| [Apache Maven GPG Plugin][41]                           | [Apache License, Version 2.0][34]              |
| [Apache Maven Source Plugin][42]                        | [Apache License, Version 2.0][34]              |
| [Apache Maven Javadoc Plugin][49]                       | [Apache License, Version 2.0][34]              |
| [Nexus Staging Maven Plugin][50]                        | [Eclipse Public License][51]                   |
| [Maven Failsafe Plugin][76]                             | [Apache License, Version 2.0][34]              |
| [JaCoCo :: Maven Plugin][52]                            | [Eclipse Public License 2.0][53]               |
| [error-code-crawler-maven-plugin][54]                   | [MIT License][55]                              |
| [Reproducible Build Maven Plugin][56]                   | [Apache 2.0][27]                               |
| [Project keeper maven plugin][43]                       | [The MIT License][44]                          |
| [Maven Clean Plugin][57]                                | [The Apache Software License, Version 2.0][27] |
| [Maven Resources Plugin][58]                            | [The Apache Software License, Version 2.0][27] |
| [Maven JAR Plugin][59]                                  | [The Apache Software License, Version 2.0][27] |
| [Maven Install Plugin][60]                              | [The Apache Software License, Version 2.0][27] |
| [Maven Site Plugin 3][61]                               | [The Apache Software License, Version 2.0][27] |

[0]: https://github.com/dop251/goja/blob/c4d370b87b45/LICENSE
[1]: https://github.com/dop251/goja_nodejs/blob/678b33ca5009/LICENSE
[2]: https://github.com/exasol/exasol-driver-go/blob/v0.4.5/LICENSE
[3]: https://github.com/exasol/exasol-test-setup-abstraction-server/blob/go-client/v0.2.4/go-client/LICENSE
[4]: https://github.com/go-chi/chi/blob/v5.0.7/LICENSE
[5]: https://github.com/sirupsen/logrus/blob/v1.9.0/LICENSE
[6]: https://github.com/stretchr/testify/blob/v1.8.0/LICENSE
[7]: https://github.com/swaggo/http-swagger/blob/v1.3.3/LICENSE
[8]: https://github.com/DATA-DOG/go-sqlmock/blob/master/LICENSE
[9]: https://github.com/Nightapes/go-rest/blob/v0.2.1/LICENSE
[10]: https://github.com/kinbiko/jsonassert/blob/HEAD/LICENSE
[11]: https://github.com/swagger-api/swagger-core/modules/swagger-annotations
[12]: http://www.apache.org/licenses/LICENSE-2.0.html
[13]: https://projects.eclipse.org/projects/ee4j.jersey/jersey-client
[14]: http://www.eclipse.org/legal/epl-2.0
[15]: https://www.gnu.org/software/classpath/license.html
[16]: http://www.eclipse.org/org/documents/edl-v10.php
[17]: https://opensource.org/licenses/BSD-2-Clause
[18]: https://creativecommons.org/publicdomain/zero/1.0/
[19]: https://asm.ow2.io/license.html
[20]: jquery.org/license
[21]: http://www.opensource.org/licenses/mit-license.php
[22]: https://www.w3.org/Consortium/Legal/copyright-documents-19990405
[23]: https://projects.eclipse.org/projects/ee4j.jersey/project/jersey-media-multipart
[24]: https://projects.eclipse.org/projects/ee4j.jersey/project/jersey-media-json-jackson
[25]: https://projects.eclipse.org/projects/ee4j.jersey/project/jersey-hk2
[26]: https://github.com/FasterXML/jackson-core
[27]: http://www.apache.org/licenses/LICENSE-2.0.txt
[28]: http://github.com/FasterXML/jackson
[29]: http://sourceforge.net/projects/migbase64/
[30]: http://en.wikipedia.org/wiki/BSD_licenses
[31]: http://sonarsource.github.io/sonar-scanner-maven/
[32]: http://www.gnu.org/licenses/lgpl.txt
[33]: https://maven.apache.org/plugins/maven-compiler-plugin/
[34]: https://www.apache.org/licenses/LICENSE-2.0.txt
[35]: https://maven.apache.org/enforcer/maven-enforcer-plugin/
[36]: https://www.mojohaus.org/flatten-maven-plugin/
[37]: https://sonatype.github.io/ossindex-maven/maven-plugin/
[38]: https://maven.apache.org/surefire/maven-surefire-plugin/
[39]: http://www.mojohaus.org/versions-maven-plugin/
[40]: https://maven.apache.org/plugins/maven-deploy-plugin/
[41]: https://maven.apache.org/plugins/maven-gpg-plugin/
[42]: https://maven.apache.org/plugins/maven-source-plugin/
[43]: https://github.com/exasol/project-keeper/
[44]: https://github.com/exasol/project-keeper/blob/main/LICENSE
[45]: http://www.mojohaus.org/exec-maven-plugin
[46]: https://github.com/swagger-api/swagger-codegen/modules/swagger-codegen-maven-plugin
[47]: http://www.mojohaus.org/build-helper-maven-plugin/
[48]: https://opensource.org/licenses/mit-license.php
[49]: https://maven.apache.org/plugins/maven-javadoc-plugin/
[50]: http://www.sonatype.com/public-parent/nexus-maven-plugins/nexus-staging/nexus-staging-maven-plugin/
[51]: http://www.eclipse.org/legal/epl-v10.html
[52]: https://www.jacoco.org/jacoco/trunk/doc/maven.html
[53]: https://www.eclipse.org/legal/epl-2.0/
[54]: https://github.com/exasol/error-code-crawler-maven-plugin/
[55]: https://github.com/exasol/error-code-crawler-maven-plugin/blob/main/LICENSE
[56]: http://zlika.github.io/reproducible-build-maven-plugin
[57]: http://maven.apache.org/plugins/maven-clean-plugin/
[58]: http://maven.apache.org/plugins/maven-resources-plugin/
[59]: http://maven.apache.org/plugins/maven-jar-plugin/
[60]: http://maven.apache.org/plugins/maven-install-plugin/
[61]: http://maven.apache.org/plugins/maven-site-plugin/
[62]: https://github.com/exasol/extension-manager/
[63]: https://github.com/exasol/extension-manager/blob/main/LICENSE
[64]: https://github.com/exasol/exasol-test-setup-abstraction-java/
[65]: https://github.com/exasol/exasol-test-setup-abstraction-java/blob/main/LICENSE
[66]: https://github.com/exasol/test-db-builder-java/
[67]: https://github.com/exasol/test-db-builder-java/blob/main/LICENSE
[68]: https://github.com/exasol/hamcrest-resultset-matcher/
[69]: https://github.com/exasol/hamcrest-resultset-matcher/blob/main/LICENSE
[70]: https://junit.org/junit5/
[71]: https://www.eclipse.org/legal/epl-v20.html
[72]: https://github.com/mockito/mockito
[73]: https://github.com/mockito/mockito/blob/main/LICENSE
[74]: https://github.com/exasol/udf-debugging-java/
[75]: https://opensource.org/licenses/MIT
[76]: https://maven.apache.org/surefire/maven-failsafe-plugin/
