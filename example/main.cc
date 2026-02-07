#include <cstddef>
#include <cstdio>
#include <fstream>
#include <iostream>
#include <string>
#include <vector>

#include "tbl_copy.h"
#include "tbl_item.h"
#include "tbl_matchmaking.h"
#include "tbl_npc.h"
#include "tbl_skill_level.h"

using namespace server::table;

static std::string getTableFileContent(const std::string &file_path)
{
    std::ifstream fs(file_path.c_str(), std::ios::binary | std::ios::ate);
    if (fs.is_open() == false) {
        ::fprintf(stderr, "can not open file %s\n",
            file_path.c_str());
        return "";
    }

    std::vector<char> input_buffer(fs.tellg());
    fs.seekg(0);
    input_buffer.assign((std::istreambuf_iterator<char>(fs)),
                         std::istreambuf_iterator<char>());

    return std::string(&input_buffer[0], input_buffer.size());
}

int main(int argc, char *argv[])
{
    std::string csv_dir = ".";
    if (argc > 1) {
        csv_dir = argv[1];
    }

    TblCopy tbl_copy;
    TblItem tbl_item;
    TblMatchmaking tbl_matchmaking;
    TblNpc tbl_npc;
    TblSkillLevel tbl_skill_level;
    std::string error_info;

    if (tbl_copy.parse(
            getTableFileContent(csv_dir + "/copy.csv"),
            &error_info) == false) {
        ::fprintf(stderr, "parse %s failed: %s\n",
            "copy.csv", error_info.c_str());
        return 1;
    }
    if (tbl_item.parse(
            getTableFileContent(csv_dir + "/item.csv"),
            &error_info) == false) {
        ::fprintf(stderr, "parse %s failed: %s\n",
            "item.csv", error_info.c_str());
        return 1;
    }
    if (tbl_matchmaking.parse(
            getTableFileContent(csv_dir + "/matchmaking.csv"),
            &error_info) == false) {
        ::fprintf(stderr, "parse %s failed: %s\n",
            "matchmaking.csv", error_info.c_str());
        return 1;
    }
    if (tbl_npc.parse(
            getTableFileContent(csv_dir + "/npc.csv"),
            &error_info) == false) {
        ::fprintf(stderr, "parse %s failed: %s\n",
            "npc.csv", error_info.c_str());
        return 1;
    }
    if (tbl_skill_level.parse(
            getTableFileContent(csv_dir + "/skill_level.csv"),
            &error_info) == false) {
        ::fprintf(stderr, "parse %s failed: %s\n",
            "skill_level.csv", error_info.c_str());
        return 1;
    }

    {
        const TblMatchmaking::Row *row =
            tbl_matchmaking.getRow(3);
        if (row != NULL) {
            ::printf("tbl_matchmaking:3:max_count: %d\n",
                row->max_count);
        }
    }

    {
        const TblSkillLevel::RowSet *row_set =
            tbl_skill_level.getRowSet(100503);
        if (row_set != NULL) {
            ::printf("tbl_skill_level:100503:range_param:p1: %d\n",
                (*row_set)[0].range_param.p1);
        }
    }

    return 0;
}
