variable "db_url" {
  type    = string
  default = "postgres://postgres.jsbobmnbqgxhpnqqmdgm:q0ODCNyMbjGe0Wux@aws-0-ap-south-1.pooler.supabase.com:5432/biddaneer?sslmode=require"
}

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./migrations/loader",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/16/dev"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "local" {
  src = data.external_schema.gorm.url
  url = var.db_url
  dev = "docker://postgres/16/dev"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
