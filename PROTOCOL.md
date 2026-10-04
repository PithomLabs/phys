<derivation-package>/
    README.md
    manifest.json
    manifest_test.go
    derivation_test.go
    
    
README.md
    human orientation:
    question, framework, regime, premises, assumptions,
    conventions, derivation outline, interpretation, limitations

manifest.json
    machine-readable derivation metadata:
    stable derivation ID, framework, entry points, assumptions,
    conventions, premises, expected result, verification path,
    limitations/anomalies/falsification conditions

manifest_test.go
    enforce the package contract:
    manifest exists and parses strictly,
    required fields are present,
    README/derivation_test.go exist,
    manifest metadata agrees with implementation constants,
    no unknown manifest fields,
    no reflection/dynamic registry

derivation_test.go
    executable derivation:
    real derivation trace,
    intermediate golden states,
    anti-hardcoding firewall,
    provenance/assumption checks,
    replay verification    