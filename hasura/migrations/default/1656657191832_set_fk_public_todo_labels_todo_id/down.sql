alter table "public"."todo_labels" drop constraint "todo_labels_todo_id_fkey",
  add constraint "todo_labels_todo_id_fkey"
  foreign key ("todo_id")
  references "public"."todos"
  ("id") on update restrict on delete restrict;
