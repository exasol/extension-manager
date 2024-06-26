package com.exasol.extensionmanager.itest;

import java.util.Locale;

/**
 * Helper class to check the operating system this Java VM runs in.
 */
public class OsCheck {

    private OsCheck() {
        // Not instantiable
    }

    /**
     * Get the suffix of native executables for the current operating system.
     *
     * @return suffix of native executables
     */
    public static String getExecutableSuffix() {
        final String os = System.getProperty("os.name", "generic").toLowerCase(Locale.ENGLISH);
        if (os.indexOf("win") >= 0) {
            return ".exe";
        } else {
            return "";
        }
    }
}