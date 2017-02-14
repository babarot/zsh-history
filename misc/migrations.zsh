#!zsh

id=1
host=$(hostname)

cat <<__BEGIN__
PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE "history" ("id" integer not null primary key autoincrement, "date" varchar(255), "dir" varchar(255), "command" varchar(255), "status" integer, "host" varchar(255));
__BEGIN__

fc -t "%F %T" -ln 1 |
while read f t cmd
do
    date="$f $t"
    printf "INSERT INTO 'history' VALUES(%d,'%s','%s','%s',%d,'%s');\n" \
        "$id" "$date" "" "$cmd" 0 "$host"
    ((id++))
done

cat <<__END__
DELETE FROM sqlite_sequence;
INSERT INTO "sqlite_sequence" VALUES('history',$id);
COMMIT;
__END__
