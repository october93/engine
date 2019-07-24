# prereqs: graphviz, java10
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
java -jar $DIR/vendor/schemaspy-6.0.0-rc2.jar -t pgsql -dp $DIR/vendor/postgresql-42.2.2.jar -db engine_local -host localhost -u postgres -p postgres -o ./schemaspy -s public -noads