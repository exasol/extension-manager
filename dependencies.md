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
| [Apache Maven Clean Plugin][35]                         | [Apache-2.0][28]                  |
| [Apache Maven Install Plugin][36]                       | [Apache-2.0][28]                  |
| [Apache Maven Resources Plugin][37]                     | [Apache-2.0][28]                  |
| [Apache Maven Site Plugin][38]                          | [Apache License, Version 2.0][28] |
| [SonarQube Scanner for Maven][39]                       | [GNU LGPL 3][40]                  |
| [Apache Maven Toolchains Plugin][41]                    | [Apache-2.0][28]                  |
| [Apache Maven Compiler Plugin][42]                      | [Apache-2.0][28]                  |
| [Apache Maven Enforcer Plugin][43]                      | [Apache-2.0][28]                  |
| [Maven Flatten Plugin][44]                              | [Apache Software Licenese][28]    |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][45] | [ASL2][46]                        |
| [Maven Surefire Plugin][47]                             | [Apache-2.0][28]                  |
| [Versions Maven Plugin][48]                             | [Apache License, Version 2.0][28] |
| [duplicate-finder-maven-plugin Maven Mojo][49]          | [Apache License 2.0][13]          |
| [Apache Maven Deploy Plugin][50]                        | [Apache-2.0][28]                  |
| [Apache Maven GPG Plugin][51]                           | [Apache-2.0][28]                  |
| [Apache Maven Source Plugin][52]                        | [Apache License, Version 2.0][28] |
| [Exec Maven Plugin][53]                                 | [Apache License 2][28]            |
| [swagger-codegen (maven-plugin)][54]                    | [Apache License 2.0][13]          |
| [Build Helper Maven Plugin][55]                         | [The MIT License][56]             |
| [Apache Maven Javadoc Plugin][57]                       | [Apache-2.0][28]                  |
| [Nexus Staging Maven Plugin][58]                        | [Eclipse Public License][59]      |
| [JaCoCo :: Maven Plugin][60]                            | [EPL-2.0][61]                     |
| [Quality Summarizer Maven Plugin][62]                   | [MIT License][63]                 |
| [error-code-crawler-maven-plugin][64]                   | [MIT License][65]                 |
| [Reproducible Build Maven Plugin][66]                   | [Apache 2.0][46]                  |

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
| [mockito-junit-jupiter][75]        | [MIT][76]                         |
| [udf-debugging-java][77]           | [MIT License][78]                 |
| [Maven Project Version Getter][79] | [MIT License][80]                 |
| [SLF4J JDK14 Provider][81]         | [MIT License][22]                 |

### Plugin Dependencies

| Dependency                                              | License                           |
| ------------------------------------------------------- | --------------------------------- |
| [Apache Maven Clean Plugin][35]                         | [Apache-2.0][28]                  |
| [Apache Maven Install Plugin][36]                       | [Apache-2.0][28]                  |
| [Apache Maven Resources Plugin][37]                     | [Apache-2.0][28]                  |
| [Apache Maven Site Plugin][38]                          | [Apache License, Version 2.0][28] |
| [SonarQube Scanner for Maven][39]                       | [GNU LGPL 3][40]                  |
| [Apache Maven Toolchains Plugin][41]                    | [Apache-2.0][28]                  |
| [Apache Maven Compiler Plugin][42]                      | [Apache-2.0][28]                  |
| [Apache Maven Enforcer Plugin][43]                      | [Apache-2.0][28]                  |
| [Maven Flatten Plugin][44]                              | [Apache Software Licenese][28]    |
| [org.sonatype.ossindex.maven:ossindex-maven-plugin][45] | [ASL2][46]                        |
| [Maven Surefire Plugin][47]                             | [Apache-2.0][28]                  |
| [Versions Maven Plugin][48]                             | [Apache License, Version 2.0][28] |
| [duplicate-finder-maven-plugin Maven Mojo][49]          | [Apache License 2.0][13]          |
| [Apache Maven Deploy Plugin][50]                        | [Apache-2.0][28]                  |
| [Apache Maven GPG Plugin][51]                           | [Apache-2.0][28]                  |
| [Apache Maven Source Plugin][52]                        | [Apache License, Version 2.0][28] |
| [Apache Maven Javadoc Plugin][57]                       | [Apache-2.0][28]                  |
| [Nexus Staging Maven Plugin][58]                        | [Eclipse Public License][59]      |
| [Maven Failsafe Plugin][82]                             | [Apache-2.0][28]                  |
| [JaCoCo :: Maven Plugin][60]                            | [EPL-2.0][61]                     |
| [Quality Summarizer Maven Plugin][62]                   | [MIT License][63]                 |
| [error-code-crawler-maven-plugin][64]                   | [MIT License][65]                 |
| [Reproducible Build Maven Plugin][66]                   | [Apache 2.0][46]                  |
| [Apache Maven JAR Plugin][83]                           | [Apache-2.0][28]                  |

## Registry

### Compile Dependencies

| Dependency               | License          |
| ------------------------ | ---------------- |
| [aws-cdk-lib][84]        | [Apache-2.0][85] |
| [constructs][86]         | [Apache-2.0][87] |
| [source-map-support][88] | [MIT][89]        |

## Registry-upload

### Compile Dependencies

| Dependency                           | License          |
| ------------------------------------ | ---------------- |
| [@aws-sdk/client-cloudformation][90] | [Apache-2.0][91] |
| [@aws-sdk/client-cloudfront][92]     | [Apache-2.0][91] |
| [@aws-sdk/client-s3][93]             | [Apache-2.0][91] |
| [follow-redirects][94]               | [MIT][95]        |
| [octokit][96]                        | [MIT][97]        |

## Parametervalidator

### Compile Dependencies

| Dependency                                  | License |
| ------------------------------------------- | ------- |
| [@exasol/extension-parameter-validator][98] | MIT     |

[0]: https://github.com/Nightapes/go-rest/blob/v0.3.3/LICENSE
[1]: https://github.com/dop251/goja/blob/79f3a7efcdbd/LICENSE
[2]: https://github.com/dop251/goja_nodejs/blob/29b559befffc/LICENSE
[3]: https://github.com/exasol/exasol-driver-go/blob/v1.0.10/LICENSE
[4]: https://github.com/exasol/exasol-test-setup-abstraction-server/blob/go-client/v0.3.10/go-client/LICENSE
[5]: https://github.com/go-chi/chi/blob/v5.1.0/LICENSE
[6]: https://github.com/sirupsen/logrus/blob/v1.9.3/LICENSE
[7]: https://github.com/stretchr/testify/blob/v1.9.0/LICENSE
[8]: https://github.com/swaggo/http-swagger/blob/v1.3.4/LICENSE
[9]: https://github.com/DATA-DOG/go-sqlmock/blob/master/LICENSE
[10]: https://github.com/kinbiko/jsonassert/blob/HEAD/LICENSE
[11]: https://cs.opensource.google/go/x/mod/+/v0.22.0:LICENSE
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
[35]: https://maven.apache.org/plugins/maven-clean-plugin/
[36]: https://maven.apache.org/plugins/maven-install-plugin/
[37]: https://maven.apache.org/plugins/maven-resources-plugin/
[38]: https://maven.apache.org/plugins/maven-site-plugin/
[39]: http://sonarsource.github.io/sonar-scanner-maven/
[40]: http://www.gnu.org/licenses/lgpl.txt
[41]: https://maven.apache.org/plugins/maven-toolchains-plugin/
[42]: https://maven.apache.org/plugins/maven-compiler-plugin/
[43]: https://maven.apache.org/enforcer/maven-enforcer-plugin/
[44]: https://www.mojohaus.org/flatten-maven-plugin/
[45]: https://sonatype.github.io/ossindex-maven/maven-plugin/
[46]: http://www.apache.org/licenses/LICENSE-2.0.txt
[47]: https://maven.apache.org/surefire/maven-surefire-plugin/
[48]: https://www.mojohaus.org/versions/versions-maven-plugin/
[49]: https://basepom.github.io/duplicate-finder-maven-plugin
[50]: https://maven.apache.org/plugins/maven-deploy-plugin/
[51]: https://maven.apache.org/plugins/maven-gpg-plugin/
[52]: https://maven.apache.org/plugins/maven-source-plugin/
[53]: https://www.mojohaus.org/exec-maven-plugin
[54]: https://github.com/swagger-api/swagger-codegen/tree/master/modules/swagger-codegen-maven-plugin
[55]: https://www.mojohaus.org/build-helper-maven-plugin/
[56]: https://spdx.org/licenses/MIT.txt
[57]: https://maven.apache.org/plugins/maven-javadoc-plugin/
[58]: http://www.sonatype.com/public-parent/nexus-maven-plugins/nexus-staging/nexus-staging-maven-plugin/
[59]: http://www.eclipse.org/legal/epl-v10.html
[60]: https://www.jacoco.org/jacoco/trunk/doc/maven.html
[61]: https://www.eclipse.org/legal/epl-2.0/
[62]: https://github.com/exasol/quality-summarizer-maven-plugin/
[63]: https://github.com/exasol/quality-summarizer-maven-plugin/blob/main/LICENSE
[64]: https://github.com/exasol/error-code-crawler-maven-plugin/
[65]: https://github.com/exasol/error-code-crawler-maven-plugin/blob/main/LICENSE
[66]: http://zlika.github.io/reproducible-build-maven-plugin
[67]: https://github.com/exasol/extension-manager/
[68]: https://github.com/exasol/extension-manager/blob/main/LICENSE
[69]: https://github.com/exasol/exasol-test-setup-abstraction-java/
[70]: https://github.com/exasol/exasol-test-setup-abstraction-java/blob/main/LICENSE
[71]: https://github.com/exasol/test-db-builder-java/
[72]: https://github.com/exasol/test-db-builder-java/blob/main/LICENSE
[73]: https://github.com/exasol/hamcrest-resultset-matcher/
[74]: https://github.com/exasol/hamcrest-resultset-matcher/blob/main/LICENSE
[75]: https://github.com/mockito/mockito
[76]: https://opensource.org/licenses/MIT
[77]: https://github.com/exasol/udf-debugging-java/
[78]: https://github.com/exasol/udf-debugging-java/blob/main/LICENSE
[79]: https://github.com/exasol/maven-project-version-getter/
[80]: https://github.com/exasol/maven-project-version-getter/blob/main/LICENSE
[81]: http://www.slf4j.org
[82]: https://maven.apache.org/surefire/maven-failsafe-plugin/
[83]: https://maven.apache.org/plugins/maven-jar-plugin/
[84]: https://registry.npmjs.org/aws-cdk-lib/-/aws-cdk-lib-2.167.1.tgz
[85]: https://github.com/aws/aws-cdk
[86]: https://registry.npmjs.org/constructs/-/constructs-10.4.2.tgz
[87]: https://github.com/aws/constructs
[88]: https://registry.npmjs.org/source-map-support/-/source-map-support-0.5.21.tgz
[89]: https://github.com/evanw/node-source-map-support
[90]: https://registry.npmjs.org/@aws-sdk/client-cloudformation/-/client-cloudformation-3.695.0.tgz
[91]: https://github.com/aws/aws-sdk-js-v3
[92]: https://registry.npmjs.org/@aws-sdk/client-cloudfront/-/client-cloudfront-3.693.0.tgz
[93]: https://registry.npmjs.org/@aws-sdk/client-s3/-/client-s3-3.693.0.tgz
[94]: https://registry.npmjs.org/follow-redirects/-/follow-redirects-1.15.9.tgz
[95]: https://github.com/follow-redirects/follow-redirects
[96]: https://registry.npmjs.org/octokit/-/octokit-4.0.2.tgz
[97]: https://github.com/octokit/octokit.js
[98]: https://registry.npmjs.org/@exasol/extension-parameter-validator/-/extension-parameter-validator-0.3.0.tgz
