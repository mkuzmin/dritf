plugins {
    kotlin("jvm") version "2.1.10"
    kotlin("plugin.serialization") version "2.1.10"

    id("de.undercouch.download") version "5.6.0"
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
            kotlin.setSrcDirs(listOf("src"))
            resources.setSrcDirs(listOf("resources"))
        }

        test {
            kotlin.setSrcDirs(listOf("test"))
        }
    }

    compilerOptions {
        extraWarnings = true
    }
}

tasks.test {
    useJUnitPlatform()
    inputs.file("../dritf.yaml")
    dependsOn("getRegions")
}

val regions = kotlin.sourceSets["main"]
    .resources.srcDirs
    .single()
    .resolve("regions.txt")
    .readLines()

regions.forEach { region ->
    val schemaDir = layout.buildDirectory.dir("schema")

    tasks.register<de.undercouch.gradle.tasks.download.Download>("download-$region") {
        group = "schema"

        src("https://schema.cloudformation.${region}.amazonaws.com/CloudformationSchema.zip")
        dest(schemaDir.get().file("${region}.zip"))

        onlyIfModified(true)
        useETag(true)
        downloadTaskDir(schemaDir)
        quiet(true)
    }

    // Extract task
    tasks.register<Sync>("unzip-$region") {
        group = "schema"
        dependsOn("download-$region")

        from(zipTree(schemaDir.get().file("${region}.zip")))
        into(schemaDir.get().dir(region))
    }
}

tasks.register("getRegions") {
    group = "schema"
    regions.forEach { region ->
        dependsOn("unzip-$region")
    }
}
