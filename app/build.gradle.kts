plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.kotlin.compose)
}

android {
    namespace = "com.xqy.nutrition"
    compileSdk = 35

    configurations.all {
        resolutionStrategy {
            // 强制使用androidx库，排除android.support库
            eachDependency {
                if (requested.group == "com.android.support") {
                    useTarget("androidx.${requested.name}:${requested.name}:${requested.version}")
                }
            }
            // 排除旧版本的appcompat
            exclude(group = "com.android.support", module = "appcompat-v7")
            exclude(group = "com.android.support", module = "support-v4")
        }
    }

    defaultConfig {
        applicationId = "com.xqy.nutrition"
        minSdk = 29
        targetSdk = 35
        versionCode = 1
        versionName = "1.0"

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
    }

    buildTypes {
        release {
            isMinifyEnabled = false
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
    kotlinOptions {
        jvmTarget = "11"
    }
    buildFeatures {
        compose = true
    }
}

dependencies {
    implementation(libs.dotenv)
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.lifecycle.runtime.ktx)
    implementation(libs.androidx.activity.compose)
    implementation(platform(libs.androidx.compose.bom))
    implementation(libs.androidx.ui)
    implementation(libs.androidx.ui.graphics)
    implementation(libs.androidx.ui.tooling.preview)
    implementation(libs.androidx.material3)
    implementation(libs.androidx.navigation.runtime.ktx)
    implementation(libs.androidx.camera.core)
    implementation(libs.core)
    implementation(libs.androidx.navigation.compose)
    implementation(libs.cronet.embedded)
    implementation(libs.androidx.camera.lifecycle)
    implementation(libs.androidx.camera.view)
    
    // Compose Pager for welcome screen
    implementation("androidx.compose.foundation:foundation:1.7.2")
    implementation("com.google.accompanist:accompanist-pager:0.32.0")
    implementation("com.google.accompanist:accompanist-pager-indicators:0.32.0")
    
    // OkHttp for network requests
    implementation(libs.okhttp)
    
    // Coroutines for asynchronous programming
    implementation(libs.kotlinx.coroutines.android)
    implementation(libs.generativeai)

    implementation(libs.coze.api)
    implementation(libs.volcenginertc)
    implementation(libs.androidx.appcompat)

    testImplementation(libs.junit)
    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(platform(libs.androidx.compose.bom))
    androidTestImplementation(libs.androidx.ui.test.junit4)
    debugImplementation(libs.androidx.ui.tooling)
    debugImplementation(libs.androidx.ui.test.manifest)
}