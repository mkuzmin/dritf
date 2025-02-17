plugins {
    kotlin("jvm") version "2.1.10"
    kotlin("plugin.serialization") version "2.1.10"
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.7.3")
    implementation("com.charleskorn.kaml:kaml:0.67.0")

    implementation("ch.qos.logback:logback-classic:1.5.12")
    implementation("aws.sdk.kotlin:account:1.4.6")

    testImplementation(platform("org.junit:junit-bom:5.11.4"))
    testImplementation("org.junit.jupiter:junit-jupiter")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

kotlin {
    sourceSets {
        main {
            kotlin.srcDirs("src")
        }

        test {
            kotlin.srcDirs("test")
            resources.srcDirs("testResources")
        }
    }

    compilerOptions {
        extraWarnings = true
    }
}

tasks.test {
    useJUnitPlatform()
    inputs.file("../dritf.yaml")
}
