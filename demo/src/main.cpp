/* System and user includes */
#include <mpi.h>
#include <iostream>

using namespace std;

/* Error handling */
#define MPI_CHECK(call) if((call) != MPI_SUCCESS) { \
    error_exit("MPI failed when calling " #call); \
}

void error_exit(const char *error) {
    cerr << error << endl;
    MPI_Abort(MPI_COMM_WORLD, -1);
}

int main(int argc, char *argv[])
{   
    /* Initialize the MPI environment */
    MPI_CHECK(MPI_Init(&argc, &argv));

    /* Get the number of processes, current rank and hostname */
    int world_size, world_rank, name_len;
    char processor_name[MPI_MAX_PROCESSOR_NAME];

    MPI_CHECK(MPI_Comm_size(MPI_COMM_WORLD, &world_size));
    MPI_CHECK(MPI_Comm_rank(MPI_COMM_WORLD, &world_rank));
    MPI_CHECK(MPI_Get_processor_name(processor_name, &name_len));


	/*
	 * YOUR CODE HERE
	 */


    MPI_CHECK(MPI_Finalize());
    return EXIT_SUCCESS;
}