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

## Extension Integration Tests Library

### Plugin Dependencies

| Dependency                                              | License                                        |
| ------------------------------------------------------- | ---------------------------------------------- |
| [SonarQube Scanner for Maven][11]                       | [GNU LGPL 3][12]                               |
| [Apache Maven Compiler Plugin][13]                      | [Apache License, Version 2.0][14]              |
| [Apache Maven Enforcer Plugin][15]                      | [Apache License, Version 2.0][14]              |
| [Maven Flatten Plugin][16]                              | [Apache Software Licenese][17]                 |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][18] | [ASL2][17]                                     |
| [Maven Surefire Plugin][19]                             | [Apache License, Version 2.0][14]              |
| [Versions Maven Plugin][20]                             | [Apache License, Version 2.0][14]              |
| [Apache Maven Deploy Plugin][21]                        | [Apache License, Version 2.0][14]              |
| [Apache Maven GPG Plugin][22]                           | [Apache License, Version 2.0][14]              |
| [Apache Maven Source Plugin][23]                        | [Apache License, Version 2.0][14]              |
| [Apache Maven Javadoc Plugin][24]                       | [Apache License, Version 2.0][14]              |
| [Nexus Staging Maven Plugin][25]                        | [Eclipse Public License][26]                   |
| [Maven Failsafe Plugin][27]                             | [Apache License, Version 2.0][14]              |
| [JaCoCo :: Maven Plugin][28]                            | [Eclipse Public License 2.0][29]               |
| [error-code-crawler-maven-plugin][30]                   | [MIT License][31]                              |
| [Reproducible Build Maven Plugin][32]                   | [Apache 2.0][17]                               |
| [Project keeper maven plugin][33]                       | [The MIT License][34]                          |
| [swagger-codegen (maven-plugin)][35]                    | [Apache License 2.0][36]                       |
| [Build Helper Maven Plugin][37]                         | [The MIT License][38]                          |
| [Maven Clean Plugin][39]                                | [The Apache Software License, Version 2.0][17] |
| [Maven Resources Plugin][40]                            | [The Apache Software License, Version 2.0][17] |
| [Maven JAR Plugin][41]                                  | [The Apache Software License, Version 2.0][17] |
| [Maven Install Plugin][42]                              | [The Apache Software License, Version 2.0][17] |
| [Maven Site Plugin 3][43]                               | [The Apache Software License, Version 2.0][17] |

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
[11]: http://sonarsource.github.io/sonar-scanner-maven/
[12]: http://www.gnu.org/licenses/lgpl.txt
[13]: https://maven.apache.org/plugins/maven-compiler-plugin/
[14]: https://www.apache.org/licenses/LICENSE-2.0.txt
[15]: https://maven.apache.org/enforcer/maven-enforcer-plugin/
[16]: https://www.mojohaus.org/flatten-maven-plugin/
[17]: http://www.apache.org/licenses/LICENSE-2.0.txt
[18]: https://sonatype.github.io/ossindex-maven/maven-plugin/
[19]: https://maven.apache.org/surefire/maven-surefire-plugin/
[20]: http://www.mojohaus.org/versions-maven-plugin/
[21]: https://maven.apache.org/plugins/maven-deploy-plugin/
[22]: https://maven.apache.org/plugins/maven-gpg-plugin/
[23]: https://maven.apache.org/plugins/maven-source-plugin/
[24]: https://maven.apache.org/plugins/maven-javadoc-plugin/
[25]: http://www.sonatype.com/public-parent/nexus-maven-plugins/nexus-staging/nexus-staging-maven-plugin/
[26]: http://www.eclipse.org/legal/epl-v10.html
[27]: https://maven.apache.org/surefire/maven-failsafe-plugin/
[28]: https://www.jacoco.org/jacoco/trunk/doc/maven.html
[29]: https://www.eclipse.org/legal/epl-2.0/
[30]: https://github.com/exasol/error-code-crawler-maven-plugin/
[31]: https://github.com/exasol/error-code-crawler-maven-plugin/blob/main/LICENSE
[32]: http://zlika.github.io/reproducible-build-maven-plugin
[33]: https://github.com/exasol/project-keeper/
[34]: https://github.com/exasol/project-keeper/blob/main/LICENSE
[35]: https://github.com/swagger-api/swagger-codegen/modules/swagger-codegen-maven-plugin
[36]: http://www.apache.org/licenses/LICENSE-2.0.html
[37]: http://www.mojohaus.org/build-helper-maven-plugin/
[38]: https://opensource.org/licenses/mit-license.php
[39]: http://maven.apache.org/plugins/maven-clean-plugin/
[40]: http://maven.apache.org/plugins/maven-resources-plugin/
[41]: http://maven.apache.org/plugins/maven-jar-plugin/
[42]: http://maven.apache.org/plugins/maven-install-plugin/
[43]: http://maven.apache.org/plugins/maven-site-plugin/
