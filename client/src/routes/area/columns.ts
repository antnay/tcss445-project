import type { ColumnDef } from "@tanstack/table-core";
import { renderComponent } from "$lib/components/ui/data-table/index.js";
import DataTableDataButton from "./data-table-data-button.svelte";

export type CrimeData = {
    distance: string;
    crimeCategory: string;
    neighborhood: string | null;
    street: string;
    city: string;
    zip: string;
    latitude: string;
    longitude: string;
    date: string;
    time: string;
}

export const columns: ColumnDef<CrimeData>[] = [
    {
        accessorKey: "distance",
        header: ({ column }) =>
            renderComponent(DataTableDataButton, {
                text: "Distance",
                onclick: column.getToggleSortingHandler(),
            }),
    },
    {
        accessorKey: "crimeCategory",
        header: ({ column }) =>
            renderComponent(DataTableDataButton, {
                text: "Crime Category",
                onclick: column.getToggleSortingHandler(),
            }),
    },
    {
        accessorKey: "neighborhood",
        header: ({ column }) =>
            renderComponent(DataTableDataButton, {
                text: "Neighborhood",
                onclick: column.getToggleSortingHandler(),
            }),
    },
    {
        accessorKey: "street",
        header: ({ column }) =>
            renderComponent(DataTableDataButton, {
                text: "Street",
                onclick: column.getToggleSortingHandler(),
            }),
    },
    {
        accessorKey: "city",
        header: ({ column }) =>
            renderComponent(DataTableDataButton, {
                text: "City",
                onclick: column.getToggleSortingHandler(),
            }),
    },
    {
        accessorKey: "zip",
        header: ({ column }) =>
            renderComponent(DataTableDataButton, {
                text: "Zip Code",
                onclick: column.getToggleSortingHandler(),
            }),
    },
    {
        accessorKey: "latitude",
        header: ({ column }) =>
            renderComponent(DataTableDataButton, {
                text: "Latitude",
                onclick: column.getToggleSortingHandler(),
            }),
    },
    {
        accessorKey: "longitude",
        header: ({ column }) =>
            renderComponent(DataTableDataButton, {
                text: "Longitude",
                onclick: column.getToggleSortingHandler(),
            }),
    },
    {
        accessorKey: "date",
        header: ({ column }) =>
            renderComponent(DataTableDataButton, {
                text: "Date",
                onclick: column.getToggleSortingHandler(),
            }),
    },
    {
        accessorKey: "time",
        header: ({ column }) =>
            renderComponent(DataTableDataButton, {
                text: "Time",
                onclick: column.getToggleSortingHandler(),
            }),
    },
];    
