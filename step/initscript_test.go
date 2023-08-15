package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_renderTemplate(t *testing.T) {
	tests := []struct {
		name      string
		inventory templateInventory
		want      string
		wantErr   bool
	}{
		{
			name: "happy path",
			inventory: templateInventory{
				CacheVersion:    "1.+",
				CacheEndpoint:   "grpcs://example.com",
				AuthToken:       "example_token",
				PushEnabled:     true,
				DebugEnabled:    true,
				ValidationLevel: "error",
			},
			want: expectedInitScriptNoMetrics,
		},
		{
			name: "invalid endpoint",
			inventory: templateInventory{
				CacheVersion:  "1.0.0",
				CacheEndpoint: "",
				AuthToken:     "example_token",
			},
			wantErr: true,
		},
		{
			name: "with metrics enabled",
			inventory: templateInventory{
				CacheVersion:    "1.+",
				CacheEndpoint:   "grpcs://example.com",
				AuthToken:       "example_token",
				PushEnabled:     true,
				DebugEnabled:    true,
				ValidationLevel: "error",
				MetricsEnabled:  true,
				MetricsVersion:  "0.+",
				MetricsEndpoint: "example.services.bitrise.io",
				MetricsPort:     443,
			},
			want: expectedInitScriptWithMetrics,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderTemplate(tt.inventory)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

const expectedInitScriptNoMetrics = `initscript {
    repositories {
        mavenCentral()
        maven {
            url 'https://jitpack.io'
        }
    }

    dependencies {
        classpath 'io.bitrise.gradle:remote-cache:1.+'
    }
}

import io.bitrise.gradle.cache.BitriseBuildCache
import io.bitrise.gradle.cache.BitriseBuildCacheServiceFactory

gradle.settingsEvaluated { settings ->
    settings.buildCache {
        local {
            enabled = false
        }

        registerBuildCacheService(BitriseBuildCache.class, BitriseBuildCacheServiceFactory.class)
        remote(BitriseBuildCache.class) {
            endpoint = 'grpcs://example.com'
            authToken = 'example_token'
            enabled = true
            push = true
            debug = true
            blobValidationLevel = 'error'
        }
    }
}
`

const expectedInitScriptWithMetrics = `initscript {
    repositories {
        mavenCentral()
        maven {
            url 'https://jitpack.io'
        }
        maven {
            url "https://plugins.gradle.org/m2/"
        }
    }

    dependencies {
        classpath 'io.bitrise.gradle:remote-cache:1.+'
        classpath 'io.bitrise.gradle:gradle-analytics:0.+'
    }
}

rootProject {
    apply plugin: io.bitrise.gradle.analytics.AnalyticsPlugin

    analytics {
        ignoreErrors = false
        bitrise {
            remote {
                authorization = 'example_token'
                endpoint = 'example.services.bitrise.io'
                port = 443
            }
        }
    }

    // Configure the analytics producer task to run at the end of the build, no matter what tasks are executed
    allprojects {
        tasks.configureEach {
            if (name != "producer") {
                // The producer task is defined in the root project only, but we are in the allprojects {} block,
                // so this special syntax is needed to reference the root project task
                finalizedBy ":producer"
            }
        }
    }
}

import io.bitrise.gradle.cache.BitriseBuildCache
import io.bitrise.gradle.cache.BitriseBuildCacheServiceFactory

gradle.settingsEvaluated { settings ->
    settings.buildCache {
        local {
            enabled = false
        }

        registerBuildCacheService(BitriseBuildCache.class, BitriseBuildCacheServiceFactory.class)
        remote(BitriseBuildCache.class) {
            endpoint = 'grpcs://example.com'
            authToken = 'example_token'
            enabled = true
            push = true
            debug = true
            blobValidationLevel = 'error'
        }
    }
}
`
