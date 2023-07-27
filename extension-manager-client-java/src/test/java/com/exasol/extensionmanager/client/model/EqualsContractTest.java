package com.exasol.extensionmanager.client.model;

import org.junit.jupiter.api.Test;

import nl.jqno.equalsverifier.EqualsVerifier;

class EqualsContractTest {
    @Test
    void testEqualsContract() {
        EqualsVerifier.simple()
                .forClasses(APIError.class, CreateInstanceRequest.class, CreateInstanceResponse.class,
                        ExtensionsResponseExtension.class, ExtensionsResponse.class, ExtensionVersion.class,
                        InstallationsResponse.class, InstallationsResponseInstallation.class, Instance.class,
                        ListInstancesResponse.class, ParamDefinition.class, ParameterValue.class,
                        UpgradeExtensionResponse.class)
                .verify();
    }
}
