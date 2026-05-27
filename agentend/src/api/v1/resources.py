import platform
import shutil
import subprocess
from dataclasses import dataclass

from fastapi import APIRouter

router = APIRouter(prefix="/v1/resources", tags=["resources"])


@dataclass
class ResourceInfo:
    used: float
    total: float
    unit: str


def _get_disk_usage() -> ResourceInfo:
    usage = shutil.disk_usage("/")
    return ResourceInfo(
        used=usage.used / 1e9,
        total=usage.total / 1e9,
        unit="GB",
    )


def _parse_vm_stat_pages(output: str, prefix: str) -> int:
    for line in output.splitlines():
        if line.startswith(prefix):
            _, _, val = line.partition(":")
            val = val.strip().replace(".", "")
            try:
                return int(val)
            except ValueError:
                return 0
    return 0


def _get_memory_usage() -> ResourceInfo:
    system = platform.system()

    if system == "Darwin":
        try:
            out = subprocess.check_output(["sysctl", "-n", "hw.memsize"], text=True)
            total_bytes = float(out.strip())
        except (subprocess.CalledProcessError, ValueError):
            return ResourceInfo(used=0, total=0, unit="GB")

        total_gb = total_bytes / 1e9

        try:
            vm_out = subprocess.check_output("vm_stat", text=True)
        except subprocess.CalledProcessError:
            return ResourceInfo(used=0, total=total_gb, unit="GB")

        free_pages = _parse_vm_stat_pages(vm_out, "Pages free")
        inactive_pages = _parse_vm_stat_pages(vm_out, "Pages inactive")
        page_size = 4096.0
        free_gb = (free_pages + inactive_pages) * page_size / 1e9
        used_gb = total_gb - free_gb

        return ResourceInfo(used=used_gb, total=total_gb, unit="GB")

    if system == "Linux":
        try:
            with open("/proc/meminfo") as f:
                info = {}
                for line in f:
                    key, _, val = line.partition(":")
                    info[key.strip()] = val.strip()
            total_kb = float(info.get("MemTotal", "0").split()[0])
            available_kb = float(info.get("MemAvailable", "0").split()[0])
            total_gb = total_kb / 1e6
            used_gb = (total_kb - available_kb) / 1e6
            return ResourceInfo(used=used_gb, total=total_gb, unit="GB")
        except (ValueError, OSError):
            return ResourceInfo(used=0, total=0, unit="GB")

    return ResourceInfo(used=0, total=0, unit="GB")


@router.get("")
async def get_resources():
    disk = _get_disk_usage()
    memory = _get_memory_usage()
    return {
        "disk": {"used": disk.used, "total": disk.total, "unit": disk.unit},
        "memory": {"used": memory.used, "total": memory.total, "unit": memory.unit},
    }
