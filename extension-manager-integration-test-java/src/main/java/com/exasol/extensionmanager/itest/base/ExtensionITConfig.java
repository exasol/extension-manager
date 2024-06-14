package com.exasol.extensionmanager.itest.base;

import java.util.Objects;

/**
 * Configuration for {@link AbstractScriptExtensionIT} and {@link AbstractVirtualSchemaExtensionIT}.
 */
public class ExtensionITConfig {
    private static final String DEFAULT_VIRTUAL_SCHEMA_NAME_PARAM_NAME = "baseVirtualSchemaName";
    private final String projectName;
    private final String extensionId;
    private final String extensionName;
    private final String extensionDescription;
    private final String currentVersion;
    private final int expectedParameterCount;
    private final String previousVersion;
    private final String previousVersionJarFile;
    private final String virtualSchemaNameParameterName;

    private ExtensionITConfig(final Builder builder) {
        this.projectName = Objects.requireNonNull(builder.projectName, "projectName");
        this.extensionId = Objects.requireNonNull(builder.extensionId, "extensionId");
        this.extensionName = Objects.requireNonNull(builder.extensionName, "extensionName");
        this.extensionDescription = Objects.requireNonNull(builder.extensionDescription, "extensionDescription");
        this.currentVersion = Objects.requireNonNull(builder.currentVersion, "currentVersion");
        this.expectedParameterCount = builder.expectedParameterCount == null ? -1 : builder.expectedParameterCount;
        this.previousVersion = builder.previousVersion;
        this.previousVersionJarFile = builder.previousVersionJarFile;
        this.virtualSchemaNameParameterName = Objects.requireNonNull(builder.virtualSchemaNameParameterName,
                "virtualSchemaNameParameterName");
    }

    /**
     * Get the project name, e.g. {@code s3-document-files-virtual-schema}.
     * 
     * @return project name
     */
    public String getProjectName() {
        return projectName;
    }

    /**
     * Get the ID of this extension, e.g. {@code s3-vs-extension.js}.
     * 
     * @return ID of this extension
     */
    public String getExtensionId() {
        return extensionId;
    }

    /**
     * Get the user visible name of this extension, e.g. {@code S3 Virtual Schema}.
     * 
     * @return name of this extension
     */
    public String getExtensionName() {
        return extensionName;
    }

    /**
     * Get the user visible description of this extension, e.g. {@code Virtual Schema for document files on AWS S3}.
     * 
     * @return description of this extension
     */
    public String getExtensionDescription() {
        return extensionDescription;
    }

    /**
     * Get the current version of this extension, e.g. {@code 1.2.3}.
     * 
     * @return current version
     */
    public String getCurrentVersion() {
        return currentVersion;
    }

    /**
     * Get the total number of parameters for this extension, incl. virtual schema name.
     * 
     * @return total number of parameters
     */
    public int getExpectedParameterCount() {
        return expectedParameterCount;
    }

    /**
     * Get the previous version of this extension, e.g. {@code 1.2.2}.
     * <p>
     * This may be {@code null} if you are just creating the first version of the extension. Once you release a second
     * version, update this to return the previous version.
     * 
     * @return previous version of this extension
     */
    public String getPreviousVersion() {
        return previousVersion;
    }

    /**
     * Get the previous version's JAR file name of this extension, e.g.
     * {@code document-files-virtual-schema-dist-7.3.6-s3-1.2.3.jar}.
     * <p>
     * This may be {@code null} if you are just creating the first version of the extension. Once you release a second
     * version, update this to return the JAR file name.
     * 
     * @return previous version's JAR file
     */
    public String getPreviousVersionJarFile() {
        return previousVersionJarFile;
    }

    /**
     * Get the parameter name for the virtual schema name. Default is {@code base-vs.virtual-schema-name}.
     * 
     * @return parameter name for the virtual schema name
     */
    public String getVirtualSchemaNameParameterName() {
        return virtualSchemaNameParameterName;
    }

    /**
     * Create builder to build {@link ExtensionITConfig}.
     *
     * @return created builder
     */
    public static Builder builder() {
        return new Builder();
    }

    /**
     * Builder to build {@link ExtensionITConfig}.
     */
    public static final class Builder {
        private String projectName;
        private String extensionId;
        private String extensionName;
        private String extensionDescription;
        private String currentVersion;
        private Integer expectedParameterCount;
        private String previousVersion;
        private String previousVersionJarFile;
        private String virtualSchemaNameParameterName = DEFAULT_VIRTUAL_SCHEMA_NAME_PARAM_NAME;

        private Builder() {
            // empty by intention
        }

        /**
         * Set the project name, e.g. {@code s3-document-files-virtual-schema}.
         *
         * @param projectName field to set
         * @return {@code this} for fluent programming
         */
        public Builder projectName(final String projectName) {
            this.projectName = projectName;
            return this;
        }

        /**
         * Set the ID of this extension, e.g. {@code s3-vs-extension.js}.
         *
         * @param extensionId field to set
         * @return {@code this} for fluent programming
         */
        public Builder extensionId(final String extensionId) {
            this.extensionId = extensionId;
            return this;
        }

        /**
         * Set the user visible name of this extension, e.g. {@code S3 Virtual Schema}.
         *
         * @param extensionName field to set
         * @return {@code this} for fluent programming
         */
        public Builder extensionName(final String extensionName) {
            this.extensionName = extensionName;
            return this;
        }

        /**
         * Set the user visible description of this extension, e.g. {@code Virtual Schema for document files on AWS S3}.
         *
         * @param extensionDescription field to set
         * @return {@code this} for fluent programming
         */
        public Builder extensionDescription(final String extensionDescription) {
            this.extensionDescription = extensionDescription;
            return this;
        }

        /**
         * Set the current version of this extension, e.g. {@code 1.2.3}.
         *
         * @param currentVersion field to set
         * @return {@code this} for fluent programming
         */
        public Builder currentVersion(final String currentVersion) {
            this.currentVersion = currentVersion;
            return this;
        }

        /**
         * Set the total number of parameters for this extension, incl. virtual schema name.
         *
         * @param expectedParameterCount field to set
         * @return {@code this} for fluent programming
         */
        public Builder expectedParameterCount(final int expectedParameterCount) {
            this.expectedParameterCount = expectedParameterCount;
            return this;
        }

        /**
         * Set the previous version of this extension, e.g. {@code 1.2.2}.
         * <p>
         * This may be {@code null} if you are just creating the first version of the extension. Once you release a
         * second version, update this to return the previous version.
         *
         * @param previousVersion field to set
         * @return {@code this} for fluent programming
         */
        public Builder previousVersion(final String previousVersion) {
            this.previousVersion = previousVersion;
            return this;
        }

        /**
         * Set the previous version's JAR file name of this extension, e.g.
         * {@code document-files-virtual-schema-dist-7.3.6-s3-1.2.3.jar}.
         * <p>
         * This may be {@code null} if you are just creating the first version of the extension. Once you release a
         * second version, update this to return the JAR file name.
         *
         * @param previousVersionJarFile field to set
         * @return {@code this} for fluent programming
         */
        public Builder previousVersionJarFile(final String previousVersionJarFile) {
            this.previousVersionJarFile = previousVersionJarFile;
            return this;
        }

        /**
         * Parameter name for the virtual schema name. Default is {@code base-vs.virtual-schema-name}.
         * <p>
         * Set this only if your virtual schema extension does not extend the base extension and uses a custom parameter
         * name, e.g. {@code virtualSchemaName}.
         * 
         * @param virtualSchemaNameParameterName parameter name for the virtual schema name
         * @return {@code this} for fluent programming
         */
        public Builder virtualSchemaNameParameterName(final String virtualSchemaNameParameterName) {
            this.virtualSchemaNameParameterName = virtualSchemaNameParameterName;
            return this;
        }

        /**
         * Builder method of the builder.
         *
         * @return built class
         */
        public ExtensionITConfig build() {
            return new ExtensionITConfig(this);
        }
    }
}
