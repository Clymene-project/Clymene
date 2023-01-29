select
    total.__timestamp,
    total.__name__,
    total.cluster,
    total.job,
    total.instance,
    (total.mem_total - mem_free.memFree - (mem_cached.memCached + mem_buffer.memBuffer + mem_sr.memSr))/1024/1024/1024 as memoryUsed
from (select TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
             "__name__"                                       AS "__name__",
             "cluster"                                        AS "cluster",
             "job"                                            AS "job",
             "instance"                                       AS "instance",
             "mode"                                           AS "mode",
             avg("value")                                     AS "mem_total"
      from "druid"."clymene"
      WHERE "__name__" = 'node_memory_MemTotal_bytes'
      GROUP BY "__name__",
               "cluster",
               "job",
               "instance",
               "mode",
               TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S')) total
        left join (SELECT TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
                          avg("value")                                     AS "memFree"
                    FROM "druid"."clymene"
                    WHERE "__name__" = 'node_memory_MemFree_bytes'
                    GROUP BY TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S')) mem_free
             on total.__timestamp = mem_free.__timestamp
        left join (SELECT TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
                          avg("value")                                     AS "memCached"
                   FROM "druid"."clymene"
                   WHERE "__name__" = 'node_memory_Cached_bytes'
                   GROUP BY TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S')) mem_cached
             on total.__timestamp = mem_cached.__timestamp
        left join (SELECT TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
                          avg("value")                                     AS "memBuffer"
                   FROM "druid"."clymene"
                   WHERE "__name__" = 'node_memory_Buffers_bytes'
                   GROUP BY TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S')) mem_buffer
             on total.__timestamp = mem_buffer.__timestamp
        left join (SELECT TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
                          avg("value")                                     AS "memSr"
                   FROM "druid"."clymene"
                   WHERE "__name__" = 'node_memory_SReclaimable_bytes'
                   GROUP BY TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S')) mem_sr
             on total.__timestamp = mem_sr.__timestamp
where total.__timestamp >= '2023-01-27 00:00:00.000000'
  AND total.__timestamp < '2023-01-28 00:00:00.000000';

select
    total.__timestamp,
    total.__name__,
    total.cluster,
    total.job,
    total.instance,
    ROUND((total.mem_total - mem_free.memFree - (mem_cached.memCached + mem_buffer.memBuffer + mem_sr.memSr)) / total.mem_total * 100, 2) as memUtilization
from (select TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
             "__name__"                                       AS "__name__",
             "cluster"                                        AS "cluster",
             "job"                                            AS "job",
             "instance"                                       AS "instance",
             "mode"                                           AS "mode",
             avg("value")                                     AS "mem_total"
      from "druid"."clymene"
      WHERE "__name__" = 'node_memory_MemTotal_bytes'
      GROUP BY "__name__",
               "cluster",
               "job",
               "instance",
               "mode",
               TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S')) total
        left join (SELECT TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
                          avg("value")                                     AS "memFree"
                    FROM "druid"."clymene"
                    WHERE "__name__" = 'node_memory_MemFree_bytes'
                    GROUP BY TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S')) mem_free
             on total.__timestamp = mem_free.__timestamp
        left join (SELECT TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
                          avg("value")                                     AS "memCached"
                   FROM "druid"."clymene"
                   WHERE "__name__" = 'node_memory_Cached_bytes'
                   GROUP BY TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S')) mem_cached
             on total.__timestamp = mem_cached.__timestamp
        left join (SELECT TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
                          avg("value")                                     AS "memBuffer"
                   FROM "druid"."clymene"
                   WHERE "__name__" = 'node_memory_Buffers_bytes'
                   GROUP BY TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S')) mem_buffer
             on total.__timestamp = mem_buffer.__timestamp
        left join (SELECT TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
                          avg("value")                                     AS "memSr"
                   FROM "druid"."clymene"
                   WHERE "__name__" = 'node_memory_SReclaimable_bytes'
                   GROUP BY TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S')) mem_sr
             on total.__timestamp = mem_sr.__timestamp
where total.__timestamp >= '2023-01-27 00:00:00.000000'
  AND total.__timestamp < '2023-01-29 00:00:00.000000'


-- CPU TIME
select cpu1."cluster",
    cpu1."job",
    cpu1."instance",
    cpu1."namespace",
    cpu1."pod",
    cpu1."cpu_time" - cpu2."cpu_time" as cpuTime,
    cpu1.__timestamp
from (select TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
             "__name__"                                       AS "__name__",
             "cluster"                                        AS "cluster",
             "job"                                            AS "job",
             "instance"                                       AS "instance",
             "namespace"                                      AS "namespace",
             "pod"                                            AS "pod",
             sum("value")                                     AS "cpu_time"
      from "druid"."clymene"
      WHERE "__name__" = 'container_cpu_usage_seconds_total'
        AND "pod" !='' AND "pod"='etcd-minikube'
      GROUP BY "__name__",
          "cluster",
          "job",
          "instance",
          "namespace",
          "pod",
          TIME_FLOOR(CAST ("__time" AS TIMESTAMP), 'PT30S')) cpu1
         left join (select TIME_FLOOR(CAST("__time" AS TIMESTAMP), 'PT30S') AS "__timestamp",
                           "__name__"                                       AS "__name__",
                           "cluster"                                        AS "cluster",
                           "job"                                            AS "job",
                           "instance"                                       AS "instance",
                           "namespace"                                      AS "namespace",
                           "pod"                                            AS "pod",
                           sum("value")                                     AS "cpu_time"
                    from "druid"."clymene"
                    WHERE "__name__" = 'container_cpu_usage_seconds_total'
                      AND "pod" !='' AND "pod"='etcd-minikube'
                    GROUP BY "__name__",
                        "cluster",
                        "job",
                        "instance",
                        "namespace",
                        "pod",
                        TIME_FLOOR(CAST ("__time" AS TIMESTAMP), 'PT30S')) cpu2
                   on cpu1."__timestamp" = cpu2."__timestamp" + INTERVAL '30' SECOND
    AND cpu1."namespace"=cpu2."namespace"
    AND cpu1."pod"=cpu2."pod"
where cpu1."__timestamp" >= MILLIS_TO_TIMESTAMP(${__from})
    AND cpu1."__timestamp" < MILLIS_TO_TIMESTAMP(${__to})
