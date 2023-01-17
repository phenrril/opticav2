<?php include_once "includes/header.php";
include "../conexion.php";
$id_user = $_SESSION['idUser'];
$permiso = "reporte";
$sql = mysqli_query($conexion, "SELECT p.*, d.* FROM permisos p INNER JOIN detalle_permisos d ON p.id = d.id_permiso WHERE d.id_usuario = $id_user AND p.nombre = '$permiso'");
$existe = mysqli_fetch_all($sql);
if (empty($existe) && $id_user != 1) {
    header("Location: permisos.php");
}
?>


<div class="row">
    <div class="col-lg-12">
        <div class="form-group">
            <h4 class="text-center">Reportes</h4><br>
        </div>
    </div>
</div>

            <form id="form_mes" method="GET">
                <div class="row">
                    <div class="col-md-4">
                        <div class="form-group">
                            <label><b>Colocar Mes &nbsp&nbsp&nbsp</b></label>
                            <td colspan=3><select id="mes" name="mes">
                                        <option value="enero">Enero</option>
                                        <option value="febrero">Febrero</option>
                                        <option value="marzo">Marzo</option>
                                        <option value="abril">Abril</option>
                                        <option value="mayo">Mayp</option>
                                        <option value="junio">Junio</option>
                                        <option value="julio">Julio</option>
                                        <option value="agosto">Agosto</option>
                                        <option value="septiembre">Septiembre</option>
                                        <option value="octubre">Octubre</option>
                                        <option value="noviembre">Noviembre</option>
                                        <option value="diciembre">Diciembre</option>
                                        </select></td>&nbsp&nbsp&nbsp
                                        <label><b>Colocar Año &nbsp&nbsp&nbsp</b></label>
                                        <td colspan=3><select id="anio" name="anio">
                                        <option value="2023" selected>2023</option>
                                        <option value="2024">2024</option>
                                        <option value="2025">2025</option>
                                        <option value="2026">2026</option>
                                        <option value="2027">2027</option>
                                        <option value="2028">2028</option>
                                        <option value="2029">2029</option>
                                        <option value="2030">2030</option>
                                        </select></td>&nbsp&nbsp&nbsp

                                        <button type="submit" class="btn btn-primary">Buscar</button>
                        </div>
                    </div>
                </div>
                <br>
            </form>
            <table class="table table-striped" id="tabla_reportes">
                <thead>
                    <tr>
                        <th>ID Venta</th>
                        <th>Nombre Cliente</th>
                        <th>Total</th>
                        <th>Abonó</th>
                        <th>Restan</th>
                        <th>Fecha</th>
                    </tr>
                </thead>
                <tbody>
                    <?php
                    $conexion = mysqli_connect("localhost", "root", "", "sis_venta");
                    if (isset($_GET['from_date']) && isset($_GET['to_date'])) {
                        $from_date = $_GET['from_date'];
                        $to_date = $_GET['to_date'];
                        //$query = "SELECT * FROM ventas WHERE fecha BETWEEN '$from_date' AND '$to_date'";
                        $query = "SELECT ventas.*, cliente.nombre FROM ventas
                        JOIN cliente ON ventas.id_cliente = cliente.idcliente
                        WHERE ventas.fecha BETWEEN '$from_date' AND '$to_date'";
                        $query_run = mysqli_query($conexion, $query);
                        if (mysqli_num_rows($query_run) > 0) {
                            foreach ($query_run as $fila) {
                    ?>
                                <tr>
                                    <td><?php echo $fila['id']; ?></td>
                                    <td><?php echo $fila['nombre']; ?></td>
                                    <td><?php echo $fila['total']; ?></td>
                                    <td><?php echo $fila['abona']; ?></td>
                                    <td><?php echo $fila['resto']; ?></td>
                                    <td><?php echo $fila['fecha']; ?></td>
                                </tr>
                            <?php
                            }
                        } else {
                            ?>
                            <tr>
                                <td><?php echo "No se encontraron resultados"; ?></td>
                        <?php
                        
                        }
                    }
                        ?>
                        </tr>
                </tbody>
            </table>
        </div>

    </div>
</div>
</div>



    






<?php include_once "includes/footer.php"; ?>