# Generate mocks
echo "Generating mocks..."
mockery --all --exported --keeptree --with-expecter --inpackage --dir .
